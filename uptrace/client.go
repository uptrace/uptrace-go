package uptrace

import (
	"context"
	"fmt"
	"strings"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const dummySpanName = "__dummy__"

// client represents Uptrace client.
type client struct {
	dsn    *internal.DSN
	tracer trace.Tracer
}

func newClient(dsn *internal.DSN) *client {
	return &client{
		dsn:    dsn,
		tracer: otel.Tracer("uptrace-go"),
	}
}

// TraceURL returns the trace URL for the span.
func (c *client) TraceURL(span trace.Span) string {
	host := strings.TrimPrefix(c.dsn.Host, "api.")
	return fmt.Sprintf("%s://%s/search/%s?q=%s",
		c.dsn.Scheme, host, c.dsn.ProjectID, span.SpanContext().TraceID())
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

	span := trace.SpanFromContext(ctx)
	isRecording := span.IsRecording()
	if !isRecording {
		_, span = c.tracer.Start(ctx, dummySpanName)
	}

	span.AddEvent(
		"log",
		trace.WithAttributes(
			attribute.String("log.severity", "panic"),
			attribute.Any("log.message", val),
		),
	)

	if !isRecording {
		// Can't use `defer span.End()` because it recovers from panic too.
		span.End()
	}

	// Re-throw the panic.
	panic(val)
}
