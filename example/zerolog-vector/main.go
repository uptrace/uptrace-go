package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/rs/zerolog"
	"github.com/uptrace/uptrace-go/uptrace"
)

func main() {
	ctx := context.Background()

	file := &lumberjack.Logger{Filename: "zerolog.log"}
	logger := zerolog.New(file).With().Timestamp().Logger()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("app_or_package_name")

	// Create a root span (a trace) to measure some operation.
	ctx, main := tracer.Start(ctx, "main-operation")
	// End the span when the operation we are measuring is done.
	defer main.End()

	// The passed ctx carries the parent span (main).
	// That is how OpenTelemetry manages span relations.
	_, child1 := tracer.Start(ctx, "child1-of-main")
	child1.SetAttributes(attribute.String("key1", "value1"))
	logger.Error().
		// trace_id and span_id are needed to properly link the log message with the span.
		Str("trace_id", child1.SpanContext().TraceID().String()).
		Str("span_id", child1.SpanContext().SpanID().String()).
		Str("foo", "bar").
		Msg("message from zerolog")
	child1.End()

	_, child2 := tracer.Start(ctx, "child2-of-main")
	child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))
	child2.End()

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}
