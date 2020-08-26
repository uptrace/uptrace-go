/*
spanexp provides span exporter for OpenTelemetry.
*/
package spanexp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/upconfig"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	apitrace "go.opentelemetry.io/otel/api/trace"
	othttptrace "go.opentelemetry.io/otel/instrumentation/httptrace"
	"go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
)

// WithBatcher is like OpenTelemetry WithBatcher but it comes with recommended options:
//   - WithBatchTimeout(5 * time.Second)
//   - WithMaxQueueSize(10000)
//   - WithMaxExportBatchSize(10000)
func WithBatcher(cfg *upconfig.Config, opts ...sdktrace.BatchSpanProcessorOption) sdktrace.ProviderOption {
	return sdktrace.WithBatcher(NewExporter(cfg), baseOpts(opts)...)
}

// NewBatchSpanProcessor is like OpenTelemetry NewBatchSpanProcessor
// but it comes with recommended options:
//   - WithBatchTimeout(5 * time.Second)
//   - WithMaxQueueSize(10000)
//   - WithMaxExportBatchSize(10000)
func NewBatchSpanProcessor(
	cfg *upconfig.Config, opts ...sdktrace.BatchSpanProcessorOption,
) (*sdktrace.BatchSpanProcessor, error) {
	return sdktrace.NewBatchSpanProcessor(NewExporter(cfg), baseOpts(opts)...)
}

func baseOpts(opts []sdktrace.BatchSpanProcessorOption) []sdktrace.BatchSpanProcessorOption {
	return append([]sdktrace.BatchSpanProcessorOption{
		sdktrace.WithBatchTimeout(5 * time.Second),
		sdktrace.WithMaxQueueSize(10000),
		sdktrace.WithMaxExportBatchSize(10000),
	}, opts...)
}

type Exporter struct {
	cfg *upconfig.Config

	endpoint string
	token    string

	tracer apitrace.Tracer
}

var _ trace.SpanBatcher = (*Exporter)(nil)

func NewExporter(cfg *upconfig.Config) *Exporter {
	cfg.Init()

	e := &Exporter{
		cfg: cfg,

		tracer: global.Tracer("github.com/uptrace/uptrace-go"),
	}

	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		internal.Logger.Print(err.Error())
		cfg.Disabled = true
	} else {
		e.endpoint = fmt.Sprintf("%s://%s/api/v1/tracing/%s/spans",
			dsn.Scheme, dsn.Host, dsn.ProjectID)
		e.token = dsn.Token
	}

	return e
}

var _ trace.SpanBatcher = (*Exporter)(nil)

func (e *Exporter) Close() error {
	return nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []*trace.SpanData) {
	if e.cfg.Disabled {
		return
	}

	var span apitrace.Span

	if e.cfg.Trace {
		ctx, span = e.tracer.Start(ctx, "ExportSpans")
		defer span.End()

		span.SetAttributes(
			kv.Int("num_span", len(spans)),
		)
	} else {
		span = apitrace.NoopSpan{}
	}

	expoSpans := make([]expoSpan, len(spans))
	m := make(map[apitrace.ID]*expoTrace, len(spans)/10)

	for i, span := range spans {
		expose := &expoSpans[i]
		initExpoSpan(expose, span)

		if trace, ok := m[span.SpanContext.TraceID]; ok {
			trace.Spans = append(trace.Spans, expose)
		} else {
			m[span.SpanContext.TraceID] = &expoTrace{
				ID:    span.SpanContext.TraceID,
				Spans: []*expoSpan{expose},
			}
		}
	}

	traces := make([]*expoTrace, 0, len(m))

	for _, trace := range m {
		traces = append(traces, trace)
	}

	if err := e.send(ctx, traces); err != nil {
		span.SetStatus(codes.Internal, "")
		span.RecordError(ctx, err)
	}
}

//------------------------------------------------------------------------------

func (e *Exporter) send(ctx context.Context, traces []*expoTrace) error {
	var span apitrace.Span

	if e.cfg.Trace {
		ctx, span = e.tracer.Start(ctx, "send")
		defer span.End()
	}

	enc := internal.GetEncoder()
	defer internal.PutEncoder(enc)

	out := map[string]interface{}{
		"traces": traces,
	}

	buf, err := enc.EncodeS2(out)
	if err != nil {
		return err
	}

	if e.cfg.Trace && e.cfg.ClientTrace {
		ctx = httptrace.WithClientTrace(ctx, othttptrace.NewClientTrace(ctx))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.token)
	req.Header.Set("Content-Type", "application/msgpack")
	req.Header.Set("Content-Encoding", "s2")

	resp, err := e.cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return statusCodeError{code: resp.StatusCode}
	}

	return nil
}

//------------------------------------------------------------------------------

type expoTrace struct {
	ID    apitrace.ID `msgpack:"id"`
	Spans []*expoSpan `msgpack:"spans"`
}

type expoSpan struct {
	ID       uint64 `msgpack:"id"`
	ParentID uint64 `msgpack:"parentId"`

	Name      string           `msgpack:"name"`
	Kind      string           `msgpack:"kind"`
	StartTime int64            `msgpack:"startTime"`
	EndTime   int64            `msgpack:"endTime"`
	Attrs     internal.KVSlice `msgpack:"attrs"`

	StatusCode    uint32 `msgpack:"statusCode"`
	StatusMessage string `msgpack:"statusMessage"`

	Events   []expoEvent      `msgpack:"events"`
	Links    []expoLink       `msgpack:"links"`
	Resource internal.KVSlice `msgpack:"resource,omitempty"`

	Tracer struct {
		Name    string `msgpack:"name"`
		Version string `msgpack:"version"`
	} `msgpack:"tracer"`
}

func initExpoSpan(expose *expoSpan, span *trace.SpanData) {
	expose.ID = asUint64(span.SpanContext.SpanID)
	expose.ParentID = asUint64(span.ParentSpanID)

	expose.Name = span.Name
	expose.Kind = span.SpanKind.String()
	expose.StartTime = span.StartTime.UnixNano()
	expose.EndTime = span.EndTime.UnixNano()
	expose.Attrs = span.Attributes

	expose.StatusCode = uint32(span.StatusCode)
	expose.StatusMessage = span.StatusMessage

	if len(span.MessageEvents) > 0 {
		expose.Events = make([]expoEvent, len(span.MessageEvents))
		for i := range span.MessageEvents {
			initExpoEvent(&expose.Events[i], &span.MessageEvents[i])
		}
	}

	if len(span.Links) > 0 {
		expose.Links = make([]expoLink, len(span.Links))
		for i := range span.Links {
			initExpoLink(&expose.Links[i], &span.Links[i])
		}
	}

	if span.Resource != nil {
		expose.Resource = span.Resource.Attributes()
	}

	expose.Tracer.Name = span.InstrumentationLibrary.Name
	expose.Tracer.Version = span.InstrumentationLibrary.Version
}

type expoEvent struct {
	Name  string           `msgpack:"name"`
	Attrs internal.KVSlice `msgpack:"attrs"`
	Time  int64            `msgpack:"time"`
}

func initExpoEvent(expose *expoEvent, event *trace.Event) {
	expose.Name = event.Name
	expose.Attrs = event.Attributes
	expose.Time = event.Time.UnixNano()
}

type expoLink struct {
	TraceID apitrace.ID      `msgpack:"traceId"`
	SpanID  uint64           `msgpack:"spanId"`
	Attrs   internal.KVSlice `msgpack:"attrs"`
}

func initExpoLink(expose *expoLink, link *apitrace.Link) {
	expose.TraceID = link.SpanContext.TraceID
	expose.SpanID = asUint64(link.SpanContext.SpanID)
	expose.Attrs = link.Attributes
}

func asUint64(b [8]byte) uint64 {
	return binary.LittleEndian.Uint64(b[:])
}

//------------------------------------------------------------------------------

type statusCodeError struct {
	code int
}

func (e statusCodeError) Error() string {
	return "got status code " + strconv.Itoa(e.code) + ", wanted 200 OK"
}
