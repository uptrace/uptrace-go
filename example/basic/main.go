package main

import (
	"context"
	"errors"
	"fmt"
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

	// Use upclient to report errors when there is no active span.
	upclient.ReportError(ctx, errors.New("Hello from uptrace-go!"))

	// Create a tracer.
	tracer := global.Tracer("github.com/your/repo")

	// Start active span.
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
	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})

	return upclient
}
