package uptrace

import (
	"context"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/uptrace/uptrace-go/internal"
)

func configureTracing(ctx context.Context, client *client, cfg *config) {
	provider := cfg.tracerProvider
	if provider == nil {
		var opts []sdktrace.TracerProviderOption

		if res := cfg.newResource(); res != nil {
			opts = append(opts, sdktrace.WithResource(res))
		}
		if cfg.traceSampler != nil {
			opts = append(opts, sdktrace.WithSampler(cfg.traceSampler))
		}

		provider = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(provider)
	}

	exp, err := otlptrace.New(ctx, otlpTraceClient(client.dsn))
	if err != nil {
		internal.Logger.Printf("otlptrace.New failed: %s", err)
		return
	}

	queueSize := queueSize()
	bsp := sdktrace.NewBatchSpanProcessor(exp,
		sdktrace.WithMaxQueueSize(queueSize),
		sdktrace.WithMaxExportBatchSize(queueSize),
		sdktrace.WithBatchTimeout(10*time.Second),
	)
	provider.RegisterSpanProcessor(bsp)

	if cfg.prettyPrint {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			internal.Logger.Printf(err.Error())
		} else {
			provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}

	client.provider = provider
}

func otlpTraceClient(dsn *DSN) otlptrace.Client {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(dsn.OTLPHost()),
		otlptracegrpc.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn.String(),
		}),
		otlptracegrpc.WithCompressor(gzip.Name),
	}

	if dsn.Scheme == "https" {
		// Create credentials using system certificates.
		creds := credentials.NewClientTLSFromCert(nil, "")
		options = append(options, otlptracegrpc.WithTLSCredentials(creds))
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.NewClient(options...)
}

func queueSize() int {
	const min = 1000
	const max = 8000

	n := (runtime.GOMAXPROCS(0) / 2) * 1000
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
