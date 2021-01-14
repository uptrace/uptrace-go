/*
spanexp provides span exporter for OpenTelemetry.
*/
package spanexp

import (
	"context"
	"fmt"
	"net/http/httptrace"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/export/trace"
	apitrace "go.opentelemetry.io/otel/trace"
)

type Exporter struct {
	cfg *Config

	wg sync.WaitGroup
	rl *internal.Gate

	endpoint string
	token    string

	tracer apitrace.Tracer

	closed uint32
}

var _ trace.SpanExporter = (*Exporter)(nil)

func NewExporter(cfg *Config) (*Exporter, error) {
	cfg.Init()

	e := &Exporter{
		cfg: cfg,
		rl:  internal.NewGate(runtime.NumCPU()),

		tracer: otel.Tracer("github.com/uptrace/uptrace-go"),
	}

	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, err
	}

	e.endpoint = fmt.Sprintf("%s://%s/api/v1/tracing/%s/spans",
		dsn.Scheme, dsn.Host, dsn.ProjectID)
	e.token = dsn.Token

	return e, nil
}

func (e *Exporter) Shutdown(context.Context) error {
	if !atomic.CompareAndSwapUint32(&e.closed, 0, 1) {
		return nil
	}
	if e.cfg.Disabled {
		return nil
	}

	e.wg.Wait()
	return nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []*trace.SpanData) error {
	if e.cfg.Disabled {
		return nil
	}

	var currSpan apitrace.Span

	if e.cfg.Trace {
		ctx, currSpan = e.tracer.Start(ctx, "ExportSpans")
		defer currSpan.End()

		currSpan.SetAttributes(
			label.Int("num_span", len(spans)),
		)
	}

	outSpans := make([]Span, 0, len(spans))

	sampler := e.cfg.Sampler.Description()
	for _, span := range spans {
		outSpans = append(outSpans, Span{})
		out := &outSpans[len(outSpans)-1]

		initUptraceSpan(out, span)
		out.Sampler = sampler

		if !e.filter(out) {
			outSpans = outSpans[:len(outSpans)-1]
		}
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

		if err := e.SendSpans(ctx, outSpans); err != nil {
			internal.Logger.Printf(ctx, "send failed: %s", err)

			if currSpan != nil {
				currSpan.SetStatus(codes.Error, err.Error())
				currSpan.RecordError(err)
			}
		}
	}()

	return nil
}

func (e *Exporter) filter(span *Span) bool {
	for _, filter := range e.cfg.Filters {
		if !filter(span) {
			return false
		}
	}
	return true
}

//------------------------------------------------------------------------------

func (e *Exporter) SendSpans(ctx context.Context, spans []Span) error {
	if e.cfg.Trace {
		var span apitrace.Span
		ctx, span = e.tracer.Start(ctx, "send")
		defer span.End()
	}

	enc := internal.GetEncoder()
	defer internal.PutEncoder(enc)

	out := map[string]interface{}{
		"spans": spans,
	}

	data, err := enc.EncodeS2(out)
	if err != nil {
		return err
	}

	if e.cfg.Trace && e.cfg.ClientTrace {
		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	}

	return internal.PostWithRetry(ctx, e.cfg.HTTPClient, e.endpoint, e.token, data)
}
