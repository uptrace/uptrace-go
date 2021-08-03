package spanexp

import (
	"encoding/binary"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type KeyValueSlice = internal.KeyValueSlice

type Span struct {
	ID       uint64        `msgpack:"id"`
	ParentID uint64        `msgpack:"parentId"`
	TraceID  trace.TraceID `msgpack:"traceId"`

	Name      string `msgpack:"name"`
	Kind      string `msgpack:"kind"`
	StartTime int64  `msgpack:"startTime"`
	EndTime   int64  `msgpack:"endTime"`

	Resource KeyValueSlice `msgpack:"resource,omitempty"`
	Attrs    KeyValueSlice `msgpack:"attrs"`

	StatusCode    string `msgpack:"statusCode"`
	StatusMessage string `msgpack:"statusMessage"`

	TracerName    string `msgpack:"tracerName"`
	TracerVersion string `msgpack:"tracerVersion"`

	Events []Event `msgpack:"events"`
	Links  []Link  `msgpack:"links"`
}

func initUptraceSpan(out *Span, in sdktrace.ReadOnlySpan) {
	spanCtx := in.SpanContext()
	out.ID = asUint64(spanCtx.SpanID())
	out.ParentID = asUint64(in.Parent().SpanID())
	out.TraceID = spanCtx.TraceID()

	out.Name = in.Name()
	out.Kind = in.SpanKind().String()
	out.StartTime = in.StartTime().UnixNano()
	out.EndTime = in.EndTime().UnixNano()

	if resource := in.Resource(); resource != nil {
		out.Resource = resource.Attributes()
	}
	out.Attrs = in.Attributes()

	status := in.Status()
	out.StatusCode = statusCode(status.Code)
	out.StatusMessage = status.Description

	if events := in.Events(); len(events) > 0 {
		out.Events = make([]Event, len(events))
		for i := range events {
			initUptraceEvent(&out.Events[i], &events[i])
		}
	}

	if links := in.Links(); len(links) > 0 {
		out.Links = make([]Link, len(links))
		for i := range links {
			initUptraceLink(&out.Links[i], &links[i])
		}
	}

	lib := in.InstrumentationLibrary()
	out.TracerName = lib.Name
	out.TracerVersion = lib.Version
}

type Event struct {
	Name  string        `msgpack:"name"`
	Attrs KeyValueSlice `msgpack:"attrs"`
	Time  int64         `msgpack:"time"`
}

func initUptraceEvent(out *Event, in *sdktrace.Event) {
	out.Name = in.Name
	out.Attrs = in.Attributes
	out.Time = in.Time.UnixNano()
}

type Link struct {
	TraceID trace.TraceID `msgpack:"traceId"`
	SpanID  uint64        `msgpack:"spanId"`
	Attrs   KeyValueSlice `msgpack:"attrs"`
}

func initUptraceLink(out *Link, in *sdktrace.Link) {
	out.TraceID = in.SpanContext.TraceID()
	out.SpanID = asUint64(in.SpanContext.SpanID())
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
