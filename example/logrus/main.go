package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/uptrace/uptrace-go/extra/otellogrus"
	"github.com/uptrace/uptrace-go/uptrace"
)

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	// Add OpenTelemetry logging hook.
	logrus.AddHook(otellogrus.NewLoggingHook())

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
