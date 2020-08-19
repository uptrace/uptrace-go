package main

import (
	"context"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

func main() {
	ctx := context.Background()

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"service.name": "my-service",
			"hostname":     hostname,
		},
	})

	tracer := upclient.Tracer("github.com/uptrace/uptrace-go/example/basic")

	err := tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {
		trace.SpanFromContext(ctx).AddEvent(ctx, "Nice operation!", kv.Int("bogons", 100))

		trace.SpanFromContext(ctx).SetAttributes(kv.String("another", "yes"))

		return nil
	})
	if err != nil {
		panic(err)
	}

	select {}
}
