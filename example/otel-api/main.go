package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	upclient := uptrace.NewClient(&uptrace.Config{
		// Set DSN or UPTRACE_DSN env var.
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "v1.0.0",
	})
	defer upclient.Close()

	ctx := context.Background()
	spansExample(ctx)
	activeSpanExample(ctx)
	activateSpanManuallyExample(ctx)
}

// This example shows how to start a span and set some attributes.
func spansExample(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "main", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// Check if span is sampled and start recording.
	if span.IsRecording() {
		span.SetAttributes(
			attribute.String("key1", "value1"),
			attribute.Int("key2", 42),
		)

		span.AddEvent("log", trace.WithAttributes(
			attribute.String("log.severity", "error"),
			attribute.String("log.message", "User not found"),
			attribute.String("enduser.id", "123"),
		))

		span.RecordError(errors.New("error1"))

		span.SetStatus(codes.Error, "error description")
	}
}

// This example shows how to get/set active span from context.
func activeSpanExample(ctx context.Context) {
	ctx, main := tracer.Start(ctx, "main")
	defer main.End()

	childCtx, child := tracer.Start(ctx, "child")
	defer child.End()

	if trace.SpanFromContext(ctx) == main {
		fmt.Println("main is active")
	}

	if trace.SpanFromContext(childCtx) == child {
		fmt.Println("child is active")
	}
}

func activateSpanManuallyExample(ctx context.Context) {
	ctx, main := tracer.Start(ctx, "main")
	defer main.End()

	ctx2 := context.TODO()
	ctx2 = trace.ContextWithSpan(ctx2, main)

	if trace.SpanFromContext(ctx) == trace.SpanFromContext(ctx2) {
		fmt.Println("span is active in multiple contexts")
	}
}
