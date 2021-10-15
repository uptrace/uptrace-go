package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"github.com/uptrace/uptrace-go/uptrace"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	// Add OpenTelemetry logging hook.
	logrus.AddHook(otellogrus.NewHook())

	tracer := otel.Tracer("app_or_package_name")

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	// You must use WithContext to propagate the active span.
	logrus.WithContext(ctx).
		WithError(errors.New("hello world")).
		WithField("foo", "bar").
		Error("something failed")

	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
