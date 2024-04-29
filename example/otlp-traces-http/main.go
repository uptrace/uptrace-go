package main

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("UPTRACE_DSN")
	if dsn == "" {
		panic("UPTRACE_DSN environment variable is required")
	}
	fmt.Println("using DSN:", dsn)

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("otlp.uptrace.dev"),
		otlptracehttp.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(1000),
		sdktrace.WithMaxExportBatchSize(1000))
	// Call shutdown to flush the buffers when program exits.
	defer bsp.Shutdown(ctx)

	resource, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", "myservice"),
			attribute.String("service.version", "1.0.0"),
		))
	if err != nil {
		panic(err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)
	tracerProvider.RegisterSpanProcessor(bsp)

	// Install our tracer provider and we are done.
	otel.SetTracerProvider(tracerProvider)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	fmt.Printf("trace: https://app.uptrace.dev/traces/%s\n", span.SpanContext().TraceID())
}
