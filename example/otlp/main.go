package main

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx := context.Background()

	shutdown, err := installOTLP(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown()

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")

	_, child1 := tracer.Start(ctx, "child1")
	child1.SetAttributes(label.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2")
	child2.SetAttributes(label.Int("key2", 42), label.Float64("key3", 123.456))
	child2.End()

	span.End()
}

func installOTLP(ctx context.Context) (func(), error) {
	creds := credentials.NewClientTLSFromCert(nil, "")
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithEndpoint("otlp.uptrace.dev:4317"),
		otlpgrpc.WithTLSCredentials(creds),
		otlpgrpc.WithHeaders(map[string]string{
			"uptrace-token": os.Getenv("UPTRACE_TOKEN"),
		}),
		otlpgrpc.WithCompressor("gzip"),
	)

	// driver := otlphttp.NewDriver(
	// 	otlphttp.WithEndpoint("otlp.uptrace.dev:443"),
	// 	otlphttp.WithHeaders(map[string]string{
	// 		"uptrace-token": os.Getenv("UPTRACE_TOKEN"),
	// 	}),
	// 	otlphttp.WithCompression(otlphttp.GzipCompression),
	// )

	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(1000),
		sdktrace.WithMaxExportBatchSize(1000))

	tracerProvider := sdktrace.NewTracerProvider()
	tracerProvider.RegisterSpanProcessor(bsp)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		bsp.Shutdown(context.TODO())
	}, nil
}
