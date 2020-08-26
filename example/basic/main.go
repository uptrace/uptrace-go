package main

import (
	"context"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/label"
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

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	tracer := upclient.Tracer("github.com/uptrace/uptrace-go/example/basic")

	ctx, span := tracer.Start(ctx, "operation")

	span.AddEvent(ctx, "type1", label.Int("bogons", 100))
	span.SetAttributes(label.String("another", "yes"))

	span.End()

	panic("something went wrong")
}
