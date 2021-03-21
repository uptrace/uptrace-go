package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "1.0.0",
	})
	defer uptrace.Shutdown(ctx)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")

	_, child1 := tracer.Start(ctx, "child1")
	child1.SetAttributes(attribute.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2")
	child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))
	child2.End()

	span.End()
	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
