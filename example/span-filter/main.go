package main

import (
	"context"
	"fmt"

	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "test",
		ServiceVersion: "v1.0.0",
	}, uptrace.WithFilter(spanFilter))
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	tracer := otel.Tracer("app_or_package_name")

	// Start active span.
	ctx, span := tracer.Start(ctx, "main span")

	{
		_, span := tracer.Start(ctx, "child1")
		span.SetAttributes(attribute.String("key1", "value1"))
		span.AddEvent("event-name", trace.WithAttributes(attribute.String("foo", "bar")))
		span.End()
	}

	{
		_, span := tracer.Start(ctx, "child2")
		span.SetAttributes(attribute.String("key2", "value2"))
		span.AddEvent("event-name", trace.WithAttributes(attribute.String("foo", "baz")))
		span.End()
	}

	span.End()
	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}

func spanFilter(span *spanexp.Span) bool {
	span.Name += " [filter]"

	return true // true keeps the span
}
