package uptrace

import (
	"context"
	"os"
	"strings"
	"sync/atomic"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ConfigureOpentelemetry configures OpenTelemetry to export data to Uptrace.
// By default it:
//   - creates tracer provider;
//   - registers Uptrace span exporter;
//   - sets tracecontext + baggage composite context propagator.
//
// You can use UPTRACE_DISABLED env var to completely skip Uptrace configuration.
func ConfigureOpentelemetry(opts ...Option) {
	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		return
	}

	ctx := context.TODO()
	conf := newConfig(opts)

	if !conf.tracingEnabled && !conf.metricsEnabled {
		return
	}

	dsn, err := ParseDSN(conf.dsn)
	if err != nil {
		internal.Logger.Printf("invalid Uptrace DSN: %s (Uptrace is disabled)", err)
		return
	}

	if dsn.Token == "<token>" {
		internal.Logger.Printf("dummy Uptrace DSN detected: %q (Uptrace is disabled)", conf.dsn)
		return
	}

	if strings.HasSuffix(dsn.Host, ":14318") {
		internal.Logger.Printf("uptrace-go uses OTLP/gRPC exporter, but got host %q", dsn.Host)
	}

	client := newClient(dsn)

	configurePropagator(conf)
	if conf.tracingEnabled {
		configureTracing(ctx, client, conf)
	}
	if conf.metricsEnabled {
		configureMetrics(ctx, client, conf)
	}

	atomicClient.Store(client)
}

func configurePropagator(conf *config) {
	textMapPropagator := conf.textMapPropagator
	if textMapPropagator == nil {
		textMapPropagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	otel.SetTextMapPropagator(textMapPropagator)
}

//------------------------------------------------------------------------------

var (
	fallbackClient = newClient(&DSN{
		Scheme:   "https",
		Host:     "api.uptrace.dev",
		GRPCPort: "4317",
		Token:    "<token>",
	})
	atomicClient atomic.Value
)

func activeClient() *client {
	v := atomicClient.Load()
	if v == nil {
		return fallbackClient
	}
	return v.(*client)
}

func TraceURL(span trace.Span) string {
	return activeClient().TraceURL(span)
}

func ReportError(ctx context.Context, err error, opts ...trace.EventOption) {
	activeClient().ReportError(ctx, err, opts...)
}

func ReportPanic(ctx context.Context) {
	activeClient().ReportPanic(ctx)
}

func Shutdown(ctx context.Context) error {
	return activeClient().Shutdown(ctx)
}

func ForceFlush(ctx context.Context) error {
	return activeClient().ForceFlush(ctx)
}

func TracerProvider() *sdktrace.TracerProvider {
	return activeClient().tp
}

// SetLogger sets the logger to the given one.
func SetLogger(logger internal.ILogger) {
	internal.Logger = logger
}
