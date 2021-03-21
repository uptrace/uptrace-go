package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/uptrace"
)

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "test",
		ServiceVersion: "v1.0.0",

		BeforeSpanSend: func(span *spanexp.Span) {
			for i := range span.Attrs {
				attr := &span.Attrs[i]
				if attr.Key == "password" {
					attr.Value = attribute.StringValue("***")
					return
				}
			}
		},
	})
	defer uptrace.Shutdown(ctx)

	tracer := otel.Tracer("app_or_package_name")

	// Start active span.
	ctx, span := tracer.Start(ctx, "main span")
	span.SetAttributes(attribute.String("password", "qwerty"))
	span.End()

	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
