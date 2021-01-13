package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/extra/otellogrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()
	upclient := newUptraceClient()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	// Add OpenTelemetry logging hook.
	logrus.AddHook(otellogrus.NewLoggingHook())

	tracer := otel.Tracer("example")

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	// You must use WithContext to propagate the active span.
	logrus.WithContext(ctx).
		WithError(errors.New("hello world")).
		WithField("foo", "bar").
		Error("something failed")

	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}

func newUptraceClient() *uptrace.Client {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})

	return upclient
}
