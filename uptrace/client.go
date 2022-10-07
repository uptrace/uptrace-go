package uptrace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const dummySpanName = "__dummy__"

// client represents Uptrace client.
type client struct {
	dsn    *DSN
	tracer trace.Tracer

	tp *sdktrace.TracerProvider
	mp *metric.MeterProvider
}

func newClient(dsn *DSN) *client {
	return &client{
		dsn:    dsn,
		tracer: otel.Tracer("uptrace-go"),
	}
}

func (c *client) Shutdown(ctx context.Context) (lastErr error) {
	if c.tp != nil {
		if err := c.tp.Shutdown(ctx); err != nil {
			lastErr = err
		}
		c.tp = nil
	}
	if c.mp != nil {
		if err := c.mp.Shutdown(ctx); err != nil {
			lastErr = err
		}
		c.mp = nil
	}
	return lastErr
}

func (c *client) ForceFlush(ctx context.Context) (lastErr error) {
	if c.tp != nil {
		if err := c.tp.ForceFlush(ctx); err != nil {
			lastErr = err
		}
	}
	if c.mp != nil {
		if err := c.mp.ForceFlush(ctx); err != nil {
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
	if c.tp != nil {
		_ = c.tp.ForceFlush(ctx)
	}
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
