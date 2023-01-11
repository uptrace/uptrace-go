package main

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<token>@uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	spansExample(ctx)
	activeSpanExample(ctx)
	activateSpanManuallyExample(ctx)
}

// This example shows how to start a span and set some attributes.
func spansExample(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "operation-name", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// To avoid expensive computations, check that span is recording before setting any attributes.
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
	}

	// Record the error and update the span status.
	if err := errors.New("error1"); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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
