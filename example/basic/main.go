package main

import (
	"context"
	"errors"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

func main() {
	ctx := context.Background()
	upclient := setupUptrace()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("Hello from uptrace-go!"))

	tracer := global.Tracer("github.com/uptrace/uptrace-go/example/basic")
	ctx, span := tracer.Start(ctx, "main span")

	{
		ctx, span := tracer.Start(ctx, "child1")
		span.SetAttributes(label.String("key1", "value1"))
		span.AddEvent(ctx, "event-name", label.String("foo", "bar"))
		span.End()
	}

	{
		ctx, span := tracer.Start(ctx, "child2")
		span.SetAttributes(label.String("key2", "value2"))
		span.AddEvent(ctx, "event-name", label.String("foo", "baz"))
		span.End()
	}

	span.End()

	// This panic will be reported to Uptrace thanks to ReportPanic above.
	panic("something went wrong")
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}
