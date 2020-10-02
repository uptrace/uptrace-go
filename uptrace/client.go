package uptrace

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/upconfig"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const dummySpanName = "__dummy__"

type Config = upconfig.Config

// Client represents Uptrace client.
type Client struct {
	cfg *Config
	dsn *internal.DSN

	tracer apitrace.Tracer

	sp       sdktrace.SpanProcessor
	provider *sdktrace.TracerProvider
}

func NewClient(cfg *Config) *Client {
	cfg.Init()

	client := &Client{
		cfg: cfg,

		tracer: global.Tracer(cfg.TracerName),
	}
	client.setupTracing()

	if dsn, err := internal.ParseDSN(cfg.DSN); err == nil {
		client.dsn = dsn
	}

	return client
}

// Closes closes the client releasing associated resources.
func (c *Client) Close() error {
	runtime.Gosched()
	c.provider.UnregisterSpanProcessor(c.sp)
	return nil
}

// TraceURL returns the trace URL for the span.
func (c *Client) TraceURL(span trace.Span) string {
	host := strings.TrimPrefix(c.dsn.Host, "api.")
	return fmt.Sprintf("%s://%s/%s/search/%s",
		c.dsn.Scheme, host, c.dsn.ProjectID, span.SpanContext().TraceID)
}

// ReportError reports an error as a span event creating a dummy span if necessary.
func (c *Client) ReportError(ctx context.Context, err error, opts ...apitrace.ErrorOption) {
	span := apitrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, dummySpanName)
		defer span.End()
	}

	span.RecordError(ctx, err, opts...)
}

// ReportPanic is used with defer to report panics.
func (c *Client) ReportPanic(ctx context.Context) {
	val := recover()
	if val == nil {
		return
	}

	span := apitrace.SpanFromContext(ctx)
	isRecording := span.IsRecording()
	if !isRecording {
		ctx, span = c.tracer.Start(ctx, dummySpanName)
	}

	span.AddEvent(
		ctx,
		"log",
		label.String("log.severity", "panic"),
		label.Any("log.message", val),
	)

	if !isRecording {
		// Can't use `defer span.End()` because it recovers from panic too.
		span.End()
	}

	// Re-throw the panic.
	panic(val)
}

//------------------------------------------------------------------------------

// Tracer is a shortcut for global.Tracer that returns a named tracer.
func (c *Client) Tracer(name string) apitrace.Tracer {
	return global.Tracer(name)
}

// WithSpan is a helper that wraps the function with a span and records the returned error.
func (c *Client) WithSpan(
	ctx context.Context,
	name string,
	fn func(ctx context.Context, span trace.Span) error,
) error {
	ctx, span := c.tracer.Start(ctx, name)
	defer span.End()

	if err := fn(ctx, span); err != nil {
		span.RecordError(ctx, err)
		return err
	}
	return nil
}

func (c *Client) setupTracing() {
	if c.cfg.Disabled {
		return
	}

	kvs := make([]label.KeyValue, 0, len(c.cfg.Resource))
	for k, v := range c.cfg.Resource {
		kvs = append(kvs, label.Any(k, v))
	}

	var res *resource.Resource
	if len(kvs) > 0 {
		res = resource.New(kvs...)
	}

	c.provider = sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{
			Resource:       res,
			DefaultSampler: c.cfg.Sampler,
		}),
	)
	c.sp = spanexp.NewBatchSpanProcessor(c.cfg)
	c.provider.RegisterSpanProcessor(c.sp)
	global.SetTracerProvider(c.provider)
}
