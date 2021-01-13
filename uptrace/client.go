package uptrace

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/internal/upconfig"
	"github.com/uptrace/uptrace-go/spanexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const dummySpanName = "__dummy__"

var defaultDSN = &internal.DSN{
	ProjectID: "<project_id>",
	Token:     "<token>",

	Scheme: "https",
	Host:   "api.uptrace.dev",
}

type Config = upconfig.Config

// Client represents Uptrace client.
type Client struct {
	cfg *Config
	dsn *internal.DSN

	tracer trace.Tracer

	sp       *spanexp.Exporter
	bsp      *sdktrace.BatchSpanProcessor
	provider *sdktrace.TracerProvider
}

func NewClient(cfg *Config) *Client {
	upconfig.Init(cfg)

	client := &Client{
		cfg: cfg,

		tracer: otel.Tracer(cfg.TracerName),
	}

	client.setupTracing()
	otel.SetTextMapPropagator(cfg.TextMapPropagator)

	if dsn, err := internal.ParseDSN(cfg.DSN); err == nil {
		client.dsn = dsn
	} else {
		client.dsn = defaultDSN
	}

	return client
}

// Closes closes the client releasing associated resources.
func (c *Client) Close() error {
	runtime.Gosched()
	c.provider.UnregisterSpanProcessor(c.bsp)
	_ = c.sp.Shutdown(context.TODO())
	return nil
}

// TraceURL returns the trace URL for the span.
func (c *Client) TraceURL(span trace.Span) string {
	host := strings.TrimPrefix(c.dsn.Host, "api.")
	return fmt.Sprintf("%s://%s/%s/search?q=%s",
		c.dsn.Scheme, host, c.dsn.ProjectID, span.SpanContext().TraceID)
}

// ReportError reports an error as a span event creating a dummy span if necessary.
func (c *Client) ReportError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		_, span = c.tracer.Start(ctx, dummySpanName)
		defer span.End()
	}

	span.RecordError(err, opts...)
}

// ReportPanic is used with defer to report panics.
func (c *Client) ReportPanic(ctx context.Context) {
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
			label.String("log.severity", "panic"),
			label.Any("log.message", val),
		),
	)

	if !isRecording {
		// Can't use `defer span.End()` because it recovers from panic too.
		span.End()
	}

	// Re-throw the panic.
	panic(val)
}

//------------------------------------------------------------------------------

// TracerProvider returns a tracer provider.
func (c *Client) TracerProvider() trace.TracerProvider {
	return c.provider
}

// Tracer returns a named tracer.
func (c *Client) Tracer(name string) trace.Tracer {
	return c.provider.Tracer(name)
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
		span.RecordError(err)
		return err
	}
	return nil
}

func (c *Client) setupTracing() {
	if c.cfg.Disabled {
		return
	}

	c.provider = sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{
			Resource:       c.cfg.Resource,
			DefaultSampler: c.cfg.Sampler,
		}),
	)
	c.sp = spanexp.NewExporter(c.cfg)
	c.bsp = sdktrace.NewBatchSpanProcessor(c.sp,
		sdktrace.WithBatchTimeout(spanexp.BatchTimeout),
		sdktrace.WithMaxQueueSize(spanexp.MaxQueueSize),
		sdktrace.WithMaxExportBatchSize(spanexp.BatchSize),
	)
	c.provider.RegisterSpanProcessor(c.bsp)

	if c.cfg.PrettyPrint {
		exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
		if err == nil {
			c.provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}

	otel.SetTracerProvider(c.provider)
}
