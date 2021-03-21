package uptrace

import (
	"context"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/spanexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// SetLogger sets the logger to the given one.
func SetLogger(logger internal.ILogger) {
	internal.Logger = logger
}

func ConfigureOpentelemetry(cfg *Config) {
	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		return
	}

	if cfg.DSN == "" {
		if dsn, ok := os.LookupEnv("UPTRACE_DSN"); ok {
			cfg.DSN = dsn
		}
	}

	setupTracing(cfg)
}

func setupTracing(cfg *Config) {
	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		internal.Logger.Printf(context.TODO(), "Uptrace is disabled: %s", err)
		return
	}

	_client.Store(newClient(dsn))

	provider := cfg.TracerProvider
	if provider == nil {
		traceConfig := sdktrace.Config{
			Resource:       cfg.resource(),
			DefaultSampler: cfg.Sampler,
		}

		provider = sdktrace.NewTracerProvider(
			sdktrace.WithConfig(traceConfig),
		)
		otel.SetTracerProvider(provider)
	}

	spe, err := spanexp.NewExporter(&spanexp.Config{
		DSN:            cfg.DSN,
		Sampler:        cfg.Sampler,
		BeforeSpanSend: cfg.BeforeSpanSend,
	})
	if err != nil {
		internal.Logger.Printf(context.TODO(), "Uptrace is disabled: %s", err)
		return
	}

	queueSize := queueSize()
	bsp := sdktrace.NewBatchSpanProcessor(spe,
		sdktrace.WithMaxQueueSize(queueSize),
		sdktrace.WithMaxExportBatchSize(queueSize),
		sdktrace.WithBatchTimeout(5*time.Second),
	)
	provider.RegisterSpanProcessor(bsp)

	if cfg.PrettyPrint {
		exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			internal.Logger.Printf(context.TODO(), err.Error())
		} else {
			provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}

	textMapPropagator := cfg.TextMapPropagator
	if textMapPropagator == nil {
		textMapPropagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}
	otel.SetTextMapPropagator(textMapPropagator)
}

func queueSize() int {
	const min = 1e3
	const max = 10e3

	n := runtime.NumCPU() * 1e3
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
	fallbackDSN = &internal.DSN{
		ProjectID: "<project_id>",
		Token:     "<token>",

		Scheme: "https",
		Host:   "api.uptrace.dev",
	}

	_client        atomic.Value
	fallbackClient = newClient(fallbackDSN)
)

func activeClient() *client {
	v := _client.Load()
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
	provider := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	return provider.Shutdown(ctx)
}
