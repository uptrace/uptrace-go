package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

	// Start a root span for a websocket connection.
	ctx, span := tracer.Start(ctx, "websocket-conn", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	fmt.Printf("websocket connection: %s\n", uptrace.TraceURL(span))

	handleWebsocketRequest(ctx)
	handleWebsocketRequest(ctx)
	handleWebsocketRequest(ctx)
}

func handleWebsocketRequest(ctx context.Context) {
	// Create another span so we can end it separately from the root span.
	// The root span can stay open for hours.
	ctx, span := tracer.Start(ctx, "websocket-conn-new-request",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// Save parent context.
	parentSpan := span

	// Create a separate trace for each websocket request.
	ctx, span = tracer.Start(ctx, "websocket-request",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithNewRoot())
	defer span.End()

	fmt.Printf("websocket request: %s\n", uptrace.TraceURL(span))

	parentSpan.AddLink(trace.LinkFromContext(ctx))
	span.SetAttributes(
		attribute.String("parent_trace_id", parentSpan.SpanContext().TraceID().String()),
	)

	time.Sleep(time.Second)
}
