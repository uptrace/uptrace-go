/*
spanexp provides span exporter for OpenTelemetry.
*/
package spanexp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/export/trace"
	apitrace "go.opentelemetry.io/otel/trace"
)

type Exporter struct {
	cfg *Config

	wg sync.WaitGroup
	rl *internal.Gate

	client   internal.SimpleClient
	endpoint string

	tracer apitrace.Tracer

	closed uint32
}

var _ trace.SpanExporter = (*Exporter)(nil)

func NewExporter(cfg *Config) (*Exporter, error) {
	cfg.init()

	e := &Exporter{
		cfg: cfg,
		rl:  internal.NewGate(runtime.NumCPU()),

		tracer: otel.Tracer("uptrace-go"),
	}

	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, err
	}

	e.client.Client = cfg.HTTPClient
	e.client.Token = dsn.Token
	e.client.MaxRetries = cfg.MaxRetries

	e.endpoint = fmt.Sprintf("%s://%s/api/v1/tracing/%s/spans",
		dsn.Scheme, dsn.Host, dsn.ProjectID)

	return e, nil
}

func (e *Exporter) Shutdown(context.Context) error {
	if !atomic.CompareAndSwapUint32(&e.closed, 0, 1) {
		return nil
	}

	e.wg.Wait()
	return nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []*trace.SpanSnapshot) error {
	var currSpan apitrace.Span

	if e.cfg.Trace {
		ctx, currSpan = e.tracer.Start(ctx, "ExportSpans")
		defer currSpan.End()

		currSpan.SetAttributes(
			attribute.Int("num_span", len(spans)),
		)
	}

	outSpans := make([]Span, 0, len(spans))

	for _, span := range spans {
		outSpans = append(outSpans, Span{})
		out := &outSpans[len(outSpans)-1]

		initUptraceSpan(out, span)
		e.cfg.BeforeSpanSend(out)
	}

	if len(outSpans) == 0 {
		return nil
	}

	e.wg.Add(1)
	e.rl.Start()

	go func() {
		if currSpan != nil {
			defer currSpan.End()
		}
		defer e.rl.Done()
		defer e.wg.Done()

		out := map[string]interface{}{
			"spans": outSpans,
		}
		if e.cfg.Sampler != nil {
			out["sampler"] = e.cfg.Sampler.Description()
		}

		if err := e.SendSpans(ctx, out); err != nil {
			if err, ok := err.(*internal.StatusCodeError); ok && err.Code() == http.StatusForbidden {
				internal.Logger.Printf(ctx, "send failed: %s (DSN=%q)", err, e.cfg.DSN)
			} else {
				internal.Logger.Printf(ctx, "send failed: %s", err)
			}

			if currSpan != nil {
				currSpan.SetStatus(codes.Error, err.Error())
				currSpan.RecordError(err)
			}
		}
	}()

	return nil
}

//------------------------------------------------------------------------------

func (e *Exporter) SendSpans(ctx context.Context, out interface{}) error {
	if e.cfg.Trace {
		var span apitrace.Span
		ctx, span = e.tracer.Start(ctx, "SendSpans")
		defer span.End()
	}

	data, err := internal.EncodeMsgpack(out)
	if err != nil {
		return err
	}

	// Create a new context since then context from Otel is canceled on shutdown.
	ctx = internal.UndoContext(ctx)

	if e.cfg.Trace && e.cfg.ClientTrace {
		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	}

	return e.client.Post(ctx, e.endpoint, data)
}
