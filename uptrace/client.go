package uptrace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const dummySpanName = "__dummy__"

// client represents Uptrace client.
type client struct {
	dsn    *DSN
	tracer trace.Tracer

	provider *sdktrace.TracerProvider
	ctrl     *controller.Controller
}

func newClient(dsn *DSN) *client {
	return &client{
		dsn:    dsn,
		tracer: otel.Tracer("uptrace-go"),
	}
}

func (c *client) Shutdown(ctx context.Context) (lastErr error) {
	if c.provider != nil {
		if err := c.provider.Shutdown(ctx); err != nil {
			lastErr = err
		}
	}
	if c.ctrl != nil {
		if err := c.ctrl.Stop(ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (c *client) ForceFlush(ctx context.Context) (lastErr error) {
	if c.provider != nil {
		if err := c.provider.ForceFlush(ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// TraceURL returns the trace URL for the span.
func (c *client) TraceURL(span trace.Span) string {
	return fmt.Sprintf("%s/traces/%s", c.dsn.AppAddr(), span.SpanContext().TraceID())
}

// ReportError reports an error as a span event creating a dummy span if necessary.
func (c *client) ReportError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, dummySpanName)
		defer span.End()
	}

	span.RecordError(err, opts...)
}

// ReportPanic is used with defer to report panics.
func (c *client) ReportPanic(ctx context.Context) {
	val := recover()
	if val == nil {
		return
	}
	c.reportPanic(ctx, val)
	// Force flush since we are about to exit on panic.
	_ = ForceFlush(ctx)
	// Re-throw the panic.
	panic(val)
}

func (c *client) reportPanic(ctx context.Context, val interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, dummySpanName)
		defer span.End()
	}

	span.AddEvent(
		"log",
		trace.WithAttributes(
			attribute.String("log.severity", "panic"),
			attribute.String("log.message", fmt.Sprint(val)),
		),
	)
}
