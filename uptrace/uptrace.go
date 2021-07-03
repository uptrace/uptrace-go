package uptrace

import (
	"context"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/metricexp"
	"github.com/uptrace/uptrace-go/spanexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// SetLogger sets the logger to the given one.
func SetLogger(logger internal.ILogger) {
	internal.Logger = logger
}

// ConfigureOpentelemetry configures OpenTelemetry to export data to Uptrace.
// By default it:
//   - creates tracer provider;
//   - registers Uptrace span exporter;
//   - sets tracecontext + baggage composite context propagator.
func ConfigureOpentelemetry(cfg *Config) {
	ctx := context.TODO()

	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		return
	}

	if cfg.DSN == "" {
		if dsn, ok := os.LookupEnv("UPTRACE_DSN"); ok {
			cfg.DSN = dsn
		}
	}

	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		internal.Logger.Printf("uptrace is disabled: %s", err)
		return
	}

	client := newClient(dsn)

	configureTracing(cfg, client)
	configurePropagator(cfg)
	if !cfg.MetricsDisabled {
		configureMetrics(ctx, cfg, client)
	}

	atomicClient.Store(client)
}

func configureTracing(cfg *Config, client *client) {
	client.provider = cfg.TracerProvider
	if client.provider == nil {
		var opts []sdktrace.TracerProviderOption

		if res := cfg.newResource(); res != nil {
			opts = append(opts, sdktrace.WithResource(res))
		}
		if cfg.Sampler != nil {
			opts = append(opts, sdktrace.WithSampler(cfg.Sampler))
		}

		client.provider = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(client.provider)
	}

	spe, err := spanexp.NewExporter(&spanexp.Config{
		DSN:            cfg.DSN,
		Sampler:        cfg.Sampler,
		BeforeSpanSend: cfg.BeforeSpanSend,
	})
	if err != nil {
		internal.Logger.Printf("spanexp.NewExporter failed: %s", err)
		return
	}

	queueSize := queueSize()
	bsp := sdktrace.NewBatchSpanProcessor(spe,
		sdktrace.WithMaxQueueSize(queueSize),
		sdktrace.WithMaxExportBatchSize(queueSize),
		sdktrace.WithBatchTimeout(5*time.Second),
	)
	client.provider.RegisterSpanProcessor(bsp)

	if cfg.PrettyPrint {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			internal.Logger.Printf(err.Error())
		} else {
			client.provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}
}

func configurePropagator(cfg *Config) {
	textMapPropagator := cfg.TextMapPropagator
	if textMapPropagator == nil {
		textMapPropagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	otel.SetTextMapPropagator(textMapPropagator)
}

func configureMetrics(ctx context.Context, cfg *Config, client *client) {
	ctrl, err := metricexp.InstallNewPipeline(ctx, &metricexp.Config{
		DSN: cfg.DSN,
	}, controller.WithResource(cfg.newResource()))
	if err != nil {
		internal.Logger.Printf("metricexp.InstallNewPipeline failed: %s", err)
		return
	}
	client.ctrl = ctrl
}

func queueSize() int {
	const min = 1e3
	const max = 10e3

	n := runtime.GOMAXPROCS(0) * 1e3
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

//------------------------------------------------------------------------------
var (
	fallbackClient = newClient(&internal.DSN{
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
