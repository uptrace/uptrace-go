package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/uptrace/uptrace-go/upmetric"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	if err := setupUptrace(ctx); err != nil {
		log.Printf("setupUptrace failed: %s", err)
	}

	ctrl, err := setupUpmetric(ctx)
	if err != nil {
		log.Printf("setupUpmetric failed: %s", err)
	} else {
		defer ctrl.Stop()
	}

	// Report number of goroutines to Uptrace.
	Meter().NewInt64ValueObserver("go.num_goroutine",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			n := int64(runtime.NumGoroutine())
			result.Observe(n)
		})

	err := Tracer().WithSpan(ctx, "operation", func(ctx context.Context) error {
		trace.SpanFromContext(ctx).AddEvent(ctx, "Nice operation!", kv.Int("bogons", 100))

		trace.SpanFromContext(ctx).SetAttributes(kv.String("another", "yes"))

		return nil
	})
	if err != nil {
		panic(err)
	}

	select {}
}

func Meter() metric.MeterMust {
	return metric.Must(global.Meter("uptrace/example/basic"))
}

func Tracer() trace.Tracer {
	return global.Tracer("uptrace/example/basic")
}

func setupUptrace(ctx context.Context) error {
	exporter := uptrace.NewExporter(&uptrace.Config{
		DSN: "", // copy your project here or use UPTRACE_DSN env var
	})

	hostname, _ := os.Hostname()
	resource := resource.New(
		standard.ServiceNameKey.String("my-service"),
		standard.HostNameKey.String(hostname),
	)

	provider, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{
			Resource:       resource,
			DefaultSampler: sdktrace.AlwaysSample(),
		}),
		sdktrace.WithBatcher(exporter, sdktrace.WithMaxExportBatchSize(10000)),
	)
	if err != nil {
		return err
	}

	global.SetTraceProvider(provider)

	return nil
}

func setupUpmetric(ctx context.Context) (*push.Controller, error) {
	hostname, _ := os.Hostname()
	resource := resource.New(
		standard.ServiceNameKey.String("my-service"),
		standard.HostNameKey.String(hostname),
	)

	ctrl := upmetric.InstallNewPipeline(&upmetric.Config{
		DSN: "", // copy your project DSN here or use UPTRACE_DSN env var
	}, push.WithResource(resource))

	return ctrl, nil
}
