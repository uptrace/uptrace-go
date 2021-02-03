package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "v1.0.0",
	})

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	// Use upclient to report errors when there are no spans.
	upclient.ReportError(ctx, errors.New("Hello from uptrace-go"))

	// Create a tracer.
	tracer := otel.Tracer("github.com/your/repo")

	// Start a span.
	ctx, span := tracer.Start(ctx, "main span")

	_, child1 := tracer.Start(ctx, "child1")
	child1.SetAttributes(label.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2")
	child2.SetAttributes(label.Int("key2", 42), label.Float64("key3", 123.456))
	child2.End()

	span.End()
	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}
