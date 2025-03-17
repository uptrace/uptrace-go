package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("UPTRACE_DSN")
	if dsn == "" {
		panic("UPTRACE_DSN environment variable is required")
	}
	fmt.Println("using DSN:", dsn)

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

	shutdownTracing := configureTracing(ctx, dsn, resource)
	defer shutdownTracing()

	shutdownLogging := configureLogging(ctx, dsn, resource)
	defer shutdownLogging()

	tracer := otel.Tracer("app_or_package_name")
	logger := otelslog.NewLogger("app_or_package_name")

	ctx, main := tracer.Start(ctx, "main-operation", trace.WithSpanKind(trace.SpanKindServer))
	defer main.End()

	logger.ErrorContext(ctx, "hello world", slog.String("error", "error message"))

	fmt.Printf("trace: https://app.uptrace.dev/traces/%s\n", main.SpanContext().TraceID())
}

func configureLogging(ctx context.Context, dsn string, resource *resource.Resource) func() {
	exp, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint("api.uptrace.dev"),
		otlploghttp.WithHeaders(map[string]string{
			"uptrace-dsn": dsn,
		}),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdklog.NewBatchProcessor(exp,
		sdklog.WithMaxQueueSize(10_000),
		sdklog.WithExportMaxBatchSize(10_000),
		sdklog.WithExportInterval(10*time.Second),
		sdklog.WithExportTimeout(10*time.Second),
	)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(bsp),
		sdklog.WithResource(resource),
	)

	global.SetLoggerProvider(provider)

	return func() {
		provider.Shutdown(ctx)
	}
}

func configureTracing(ctx context.Context, dsn string, resource *resource.Resource) func() {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("api.uptrace.dev"),
		otlptracehttp.WithHeaders(map[string]string{
			"uptrace-dsn": dsn,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(10_000),
		sdktrace.WithMaxExportBatchSize(10_000),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
	)
	tracerProvider.RegisterSpanProcessor(bsp)

	otel.SetTracerProvider(tracerProvider)

	return func() {
		tracerProvider.Shutdown(ctx)
	}
}
