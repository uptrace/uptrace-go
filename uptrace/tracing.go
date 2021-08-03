package uptrace

import (
	"context"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/spanexp"
)

func configureTracing(ctx context.Context, client *client, cfg *config) {
	provider := cfg.TracerProvider
	if provider == nil {
		var opts []sdktrace.TracerProviderOption

		if res := cfg.newResource(); res != nil {
			opts = append(opts, sdktrace.WithResource(res))
		}
		if cfg.TraceSampler != nil {
			opts = append(opts, sdktrace.WithSampler(cfg.TraceSampler))
		}

		provider = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(provider)
	}

	exp, err := spanexp.NewExporter(&spanexp.Config{
		DSN:     cfg.DSN,
		Sampler: cfg.TraceSampler,
	})
	if err != nil {
		internal.Logger.Printf("spanexp.NewExporter failed: %s", err)
		return
	}

	// exp, err := otlptrace.New(ctx, otlpTraceClient(client.dsn))
	// if err != nil {
	// 	internal.Logger.Printf("otlptrace.New failed: %s", err)
	// 	return
	// }

	queueSize := queueSize()
	bsp := sdktrace.NewBatchSpanProcessor(exp,
		sdktrace.WithMaxQueueSize(queueSize),
		sdktrace.WithMaxExportBatchSize(queueSize),
		sdktrace.WithBatchTimeout(10*time.Second),
	)
	provider.RegisterSpanProcessor(bsp)

	if cfg.PrettyPrint {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			internal.Logger.Printf(err.Error())
		} else {
			provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}

	client.provider = provider
}

// func otlpTraceClient(dsn *internal.DSN) otlptrace.Client {
// 	endpoint := dsn.OTLPEndpoint()

// 	options := []otlptracegrpc.Option{
// 		otlptracegrpc.WithEndpoint(endpoint),
// 		otlptracegrpc.WithHeaders(map[string]string{
// 			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
// 			"uptrace-dsn": dsn.String(),
// 		}),
// 		otlptracegrpc.WithCompressor(gzip.Name),
// 	}

// 	if dsn.Scheme == "https" {
// 		// Create credentials using system certificates.
// 		creds := credentials.NewClientTLSFromCert(nil, "")
// 		options = append(options, otlptracegrpc.WithTLSCredentials(creds))
// 	} else {
// 		options = append(options, otlptracegrpc.WithInsecure())
// 	}

// 	return otlptracegrpc.NewClient(options...)
// }

func queueSize() int {
	const min = 1000
	const max = 8000

	n := runtime.GOMAXPROCS(0) * 1000
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
