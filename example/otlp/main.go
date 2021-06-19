package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

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

	tracerProvider := sdktrace.NewTracerProvider()
	tracerProvider.RegisterSpanProcessor(bsp)

	// Install our tracer provider and we are done.
	otel.SetTracerProvider(tracerProvider)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	_, child1 := tracer.Start(ctx, "child1")
	child1.SetAttributes(attribute.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2")
	child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))
	child2.End()

	fmt.Println("trace id:", span.SpanContext().TraceID())
}
