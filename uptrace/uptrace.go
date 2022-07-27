package uptrace

import (
	"context"
	"os"
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
	cfg := newConfig(opts)

	if !cfg.tracingEnabled && !cfg.metricsEnabled {
		return
	}

	dsn, err := ParseDSN(cfg.dsn)
	if err != nil {
		internal.Logger.Printf("uptrace is disabled: %s", err)
		return
	}

	client := newClient(dsn)

	configurePropagator(cfg)
	if cfg.tracingEnabled {
		configureTracing(ctx, client, cfg)
	}
	if cfg.metricsEnabled {
		configureMetrics(ctx, client, cfg)
	}

	atomicClient.Store(client)
}

func configurePropagator(cfg *config) {
	textMapPropagator := cfg.textMapPropagator
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
		ProjectID: "<project_id>",
		Token:     "<token>",

		Scheme: "https",
		Host:   "api.uptrace.dev",
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
