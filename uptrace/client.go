package uptrace

import (
	"context"
	"sync"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/upconfig"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config = upconfig.Config

const dummySpanName = "__dummy__"

type Client struct {
	cfg *Config

	tracer apitrace.Tracer

	setupTracingOnce sync.Once
	sp               sdktrace.SpanProcessor
	provider         *sdktrace.Provider
}

func NewClient(cfg *Config) *Client {
	cfg.Init()

	return &Client{
		cfg: cfg,

		tracer: global.Tracer("github.com/uptrace/uptrace-go"),
	}
}

// Closes closes the client releasing associated resources.
func (c *Client) Close() error {
	c.provider.UnregisterSpanProcessor(c.sp)
	return nil
}

// ReportError reports an error as a span event creating a dummy span if necessary.
func (c *Client) ReportError(ctx context.Context, err error, opts ...apitrace.ErrorOption) {
	c.setupTracing()

	span := apitrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		ctx, span = c.tracer.Start(ctx, dummySpanName)
		defer span.End()
	}

	span.RecordError(ctx, err, opts...)
}

//------------------------------------------------------------------------------

// Tracer returns a named tracer that exports span to Uptrace.
func (c *Client) Tracer(name string) apitrace.Tracer {
	c.setupTracing()

	return global.Tracer(name)
}

func (c *Client) setupTracing() {
	c.setupTracingOnce.Do(c._setupTracing)
}

func (c *Client) _setupTracing() {
	if c.cfg.Disabled {
		return
	}

	var err error

	kvs := make([]kv.KeyValue, 0, len(c.cfg.Resource))
	for k, v := range c.cfg.Resource {
		kvs = append(kvs, kv.Any(k, v))
	}

	var res *resource.Resource
	if len(kvs) > 0 {
		res = resource.New(kvs...)
	}

	c.provider, err = sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{
			Resource:       res,
			DefaultSampler: c.cfg.Sampler,
		}),
	)
	if err != nil {
		internal.Logger.Printf("NewProvider failed: %s", err)
		return
	}

	c.sp, err = spanexp.NewBatchSpanProcessor(c.cfg)
	if err != nil {
		internal.Logger.Printf("NewBatchSpanProcessor failed: %s", err)
		return
	}

	c.provider.RegisterSpanProcessor(c.sp)
	global.SetTraceProvider(c.provider)
}
