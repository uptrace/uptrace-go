package uptrace

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/uptrace/uptrace-go/internal"
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

// Client represents Uptrace client.
type Client struct {
	cfg *Config
	dsn *internal.DSN

	tracer trace.Tracer

	provider *sdktrace.TracerProvider
	bsp      *sdktrace.BatchSpanProcessor
}

func NewClient(cfg *Config, opts ...Option) *Client {
	cfg.Init(opts...)

	client := &Client{
		cfg: cfg,

		tracer: otel.Tracer(cfg.TracerName),
	}

	client.setupTracing()

	if dsn, err := internal.ParseDSN(cfg.DSN); err == nil {
		client.dsn = dsn
	} else {
		client.dsn = defaultDSN
	}

	return client
}

// Closes closes the client releasing associated resources.
func (c *Client) Close() error {
	if c.provider != nil {
		runtime.Gosched()
		c.provider.UnregisterSpanProcessor(c.bsp)
	}
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
	return c.cfg.TracerProvider
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

	if c.cfg.TracerProvider == nil {
		c.provider = c.getTracerProvider()
		c.cfg.TracerProvider = c.provider
	}

	otel.SetTextMapPropagator(c.cfg.TextMapPropagator)
	otel.SetTracerProvider(c.cfg.TracerProvider)
}

func (c *Client) getTracerProvider() *sdktrace.TracerProvider {
	const batchTimeout = 5 * time.Second

	traceConfig := sdktrace.Config{
		Resource: c.cfg.Resource,
	}
	if c.cfg.Sampler != nil {
		traceConfig.DefaultSampler = c.cfg.Sampler
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(traceConfig),
	)

	spe, err := spanexp.NewExporter(c.cfg)
	if err != nil {
		internal.Logger.Printf(context.TODO(),
			"Uptrace is disabled: %s",
			strings.TrimPrefix(err.Error(), "uptrace: "))
	} else {
		queueSize := queueSize()
		c.bsp = sdktrace.NewBatchSpanProcessor(spe,
			sdktrace.WithMaxQueueSize(queueSize),
			sdktrace.WithMaxExportBatchSize(queueSize),
			sdktrace.WithBatchTimeout(batchTimeout),
		)
		provider.RegisterSpanProcessor(c.bsp)
	}

	if c.cfg.PrettyPrint {
		exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			internal.Logger.Printf(context.TODO(), err.Error())
		} else {
			provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
		}
	}

	return provider
}

func queueSize() int {
	const min = 1e3
	const max = 10e3

	n := runtime.NumCPU() * 2e3
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
