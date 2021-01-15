package spanexp

import (
	"encoding/binary"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/export/trace"
	apitrace "go.opentelemetry.io/otel/trace"
)

type KeyValueSlice = internal.KeyValueSlice

type Span struct {
	ID       uint64           `msgpack:"id"`
	ParentID uint64           `msgpack:"parentId"`
	TraceID  apitrace.TraceID `msgpack:"traceId"`

	Name      string `msgpack:"name"`
	Kind      string `msgpack:"kind"`
	StartTime int64  `msgpack:"startTime"`
	EndTime   int64  `msgpack:"endTime"`

	Resource KeyValueSlice `msgpack:"resource,omitempty"`
	Attrs    KeyValueSlice `msgpack:"attrs"`

	StatusCode    string `msgpack:"statusCode"`
	StatusMessage string `msgpack:"statusMessage"`

	Events []Event `msgpack:"events"`
	Links  []Link  `msgpack:"links"`

	Tracer struct {
		Name    string `msgpack:"name"`
		Version string `msgpack:"version"`
	} `msgpack:"tracer"`

	Sampler string `msgpack:"sampler"`
}

func initUptraceSpan(out *Span, in *trace.SpanSnapshot) {
	out.ID = asUint64(in.SpanContext.SpanID)
	out.ParentID = asUint64(in.ParentSpanID)
	out.TraceID = in.SpanContext.TraceID

	out.Name = in.Name
	out.Kind = in.SpanKind.String()
	out.StartTime = in.StartTime.UnixNano()
	out.EndTime = in.EndTime.UnixNano()

	if in.Resource != nil {
		out.Resource = in.Resource.Attributes()
	}
	out.Attrs = in.Attributes

	out.StatusCode = statusCode(in.StatusCode)
	out.StatusMessage = in.StatusMessage

	if len(in.MessageEvents) > 0 {
		out.Events = make([]Event, len(in.MessageEvents))
		for i := range in.MessageEvents {
			initUptraceEvent(&out.Events[i], &in.MessageEvents[i])
		}
	}

	if len(in.Links) > 0 {
		out.Links = make([]Link, len(in.Links))
		for i := range in.Links {
			initUptraceLink(&out.Links[i], &in.Links[i])
		}
	}

	out.Tracer.Name = in.InstrumentationLibrary.Name
	out.Tracer.Version = in.InstrumentationLibrary.Version
}

type Event struct {
	Name  string        `msgpack:"name"`
	Attrs KeyValueSlice `msgpack:"attrs"`
	Time  int64         `msgpack:"time"`
}

func initUptraceEvent(out *Event, in *trace.Event) {
	out.Name = in.Name
	out.Attrs = in.Attributes
	out.Time = in.Time.UnixNano()
}

type Link struct {
	TraceID apitrace.TraceID `msgpack:"traceId"`
	SpanID  uint64           `msgpack:"spanId"`
	Attrs   KeyValueSlice    `msgpack:"attrs"`
}

func initUptraceLink(out *Link, in *apitrace.Link) {
	out.TraceID = in.SpanContext.TraceID
	out.SpanID = asUint64(in.SpanContext.SpanID)
	out.Attrs = in.Attributes
}

func asUint64(b [8]byte) uint64 {
	return binary.LittleEndian.Uint64(b[:])
}

func statusCode(code codes.Code) string {
	switch code {
	case codes.Unset:
		return "unset"
	case codes.Ok:
		return "ok"
	case codes.Error:
		return "error"
	default:
		return "unset"
	}
}
