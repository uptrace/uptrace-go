/*
uptrace provides span exporter for OpenTelemetry.
*/
package uptrace

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/uptrace/uptrace-go/internal"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/instrumentation/othttp"
	"go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
)

// Config is the configuration to be used when initializing an exporter.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// Disabled disables the exporter.
	// The default is to use UPTRACE_DISABLED environment var.
	Disabled bool

	endpoint string
	token    string
}

func (cfg *Config) init() {
	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		cfg.Disabled = true
		return
	}

	dsnStr := cfg.DSN
	if dsnStr == "" {
		dsnStr = os.Getenv("UPTRACE_DSN")
	}

	dsn, err := internal.ParseDSN(dsnStr)
	if err != nil {
		internal.Logger.Print(err.Error())
		cfg.Disabled = true
		return
	}

	cfg.endpoint = fmt.Sprintf("%s://%s/api/v1/tracing/%s/spans",
		dsn.Scheme, dsn.Host, dsn.ProjectID)
	cfg.token = dsn.Token
}

type Exporter struct {
	cfg *Config

	client *http.Client
}

func NewExporter(cfg *Config) *Exporter {
	cfg.init()

	client := &http.Client{
		Transport: othttp.NewTransport(http.DefaultTransport),
	}

	e := &Exporter{
		cfg: cfg,

		client: client,
	}
	return e
}

func WithBatcher(cfg *Config) sdktrace.ProviderOption {
	return sdktrace.WithBatcher(NewExporter(cfg),
		sdktrace.WithBatchTimeout(5*time.Second),
		sdktrace.WithMaxQueueSize(10000),
		sdktrace.WithMaxExportBatchSize(10000))
}

var _ trace.SpanBatcher = (*Exporter)(nil)

func (e *Exporter) Close() error {
	return nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []*trace.SpanData) {
	if e.cfg.Disabled {
		return
	}

	tracer := global.Tracer("github.com/uptrace/uptrace-go")

	ctx, span := tracer.Start(ctx, "ExportSpans")
	defer span.End()

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

	span.SetAttributes(
		kv.Int("num_span", len(spans)),
		kv.Int("num_trace", len(traces)),
	)

	_ = tracer.WithSpan(ctx, "send", func(ctx context.Context) error {
		if err := e.send(ctx, traces); err != nil {
			span.SetStatus(codes.Internal, "")
			span.AddEvent(ctx, "error",
				kv.String("error.type", reflect.TypeOf(err).String()),
				kv.String("error.message", err.Error()),
			)
			return err
		}
		return nil
	})
}

//------------------------------------------------------------------------------

func (e *Exporter) send(ctx context.Context, traces []*expoTrace) error {
	enc := internal.GetEncoder()
	defer internal.PutEncoder(enc)

	out := map[string]interface{}{
		"traces": traces,
	}

	buf, err := enc.EncodeS2(out)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.cfg.endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.cfg.token)
	req.Header.Set("Content-Type", "application/msgpack")
	req.Header.Set("Content-Encoding", "s2")

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return statusCodeError{
			code: resp.StatusCode,
		}
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

	Name          string           `msgpack:"name"`
	Kind          string           `msgpack:"kind"`
	StartTime     int64            `msgpack:"startTime"`
	EndTime       int64            `msgpack:"endTime"`
	StatusCode    uint32           `msgpack:"statusCode"`
	StatusMessage string           `msgpack:"statusMessage"`
	Attrs         internal.KVSlice `msgpack:"attrs"`

	Events   []expoEvent      `msgpack:"events"`
	Links    []expoLink       `msgpack:"links"`
	Resource internal.KVSlice `msgpack:"resource,omitempty"`
}

func initExpoSpan(expose *expoSpan, span *trace.SpanData) {
	expose.ID = asUint64(span.SpanContext.SpanID)
	expose.ParentID = asUint64(span.ParentSpanID)

	expose.Name = span.Name
	expose.Kind = span.SpanKind.String()
	expose.StartTime = span.StartTime.UnixNano()
	expose.EndTime = span.EndTime.UnixNano()
	expose.StatusCode = uint32(span.StatusCode)
	expose.StatusMessage = span.StatusMessage
	expose.Attrs = span.Attributes

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
