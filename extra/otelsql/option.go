package otelsql

import (
	"context"
	"database/sql/driver"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	provider trace.TracerProvider
	tracer   trace.Tracer
	attrs    []attribute.KeyValue
}

func newConfig(opts ...Option) *config {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	if c.provider == nil {
		c.provider = otel.GetTracerProvider()
	}
	c.tracer = c.provider.Tracer("github.com/uptrace/uptrace-go/extra/otelsql")

	return c
}

func (c *config) withSpan(
	ctx context.Context, name string, fn func(ctx context.Context, span trace.Span) error,
) error {
	ctx, span := c.tracer.Start(ctx, name, trace.WithAttributes(c.attrs...))
	err := fn(ctx, span)
	span.End()

	if !span.IsRecording() {
		return err
	}

	switch err {
	case nil,
		driver.ErrSkip,
		io.EOF: // end of rows
		// ignore
	default:
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

func (c *config) formatQuery(query string) string {
	return query
}

type Option func(c *config)

// WithTracerProvider configures a tracer provider that is used to create a tracer.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(c *config) {
		c.provider = provider
	}
}

// WithAttributes configures attributes that are used to create a span.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(c *config) {
		c.attrs = append(c.attrs, attrs...)
	}
}

// WithDBSystem configures a db.system attribute. You should prefer using
// WithAttributes and semconv, for example, `otelsql.WithAttributes(semconv.DBSystemSqlite)`.
func WithDBSystem(system string) Option {
	return func(c *config) {
		c.attrs = append(c.attrs, semconv.DBSystemKey.String(system))
	}
}
