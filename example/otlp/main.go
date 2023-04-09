package main

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx := context.Background()

	// Create credentials using system certificates.
	creds := credentials.NewClientTLSFromCert(nil, "")
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("otlp.uptrace.dev:4317"),
		otlptracegrpc.WithTLSCredentials(creds),
		otlptracegrpc.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": os.Getenv("UPTRACE_DSN"),
		}),
		otlptracegrpc.WithCompressor("gzip"),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(1000),
		sdktrace.WithMaxExportBatchSize(1000))
	// Call shutdown to flush the buffers when program exits.
	defer bsp.Shutdown(ctx)

	idg := xray.NewIDGenerator()
	tracerProvider := sdktrace.NewTracerProvider(trace.WithIDGenerator(idg))
	tracerProvider.RegisterSpanProcessor(bsp)

	// Install our tracer provider and we are done.
	otel.SetTracerProvider(tracerProvider)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	fmt.Println("trace id:", span.SpanContext().TraceID())
}
