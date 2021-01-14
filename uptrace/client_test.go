package uptrace_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/klauspost/compress/s2"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/uptrace"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

func TestFilters(t *testing.T) {
	ctx := context.Background()

	var got *spanexp.Span

	filter := func(span *spanexp.Span) bool {
		got = span
		return false
	}

	upclient := uptrace.NewClient(&uptrace.Config{
		DSN: "https://key@api.uptrace.dev/1",

		ServiceName: "test-filters",
	}, uptrace.WithFilter(filter))

	tracer := otel.Tracer("github.com/your/repo")
	_, span := tracer.Start(ctx, "main span")
	span.End()

	err := upclient.Close()
	require.Nil(t, err)

	require.NotNil(t, got)
	require.Equal(t, "main span", got.Name)

	set := label.NewSet(got.Resource...)
	val, ok := set.Value(semconv.ServiceNameKey)
	require.True(t, ok)
	require.Equal(t, "test-filters", val.AsString())
}

func TestExporter(t *testing.T) {
	ctx := context.Background()

	type In struct {
		Spans []spanexp.Span `msgpack:"spans"`
	}

	var in In

	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, "application/msgpack", req.Header.Get("Content-Type"))
		require.Equal(t, "s2", req.Header.Get("Content-Encoding"))

		b, err := ioutil.ReadAll(s2.NewReader(req.Body))
		require.NoError(t, err)

		err = msgpack.Unmarshal(b, &in)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	upclient := uptrace.NewClient(&uptrace.Config{
		DSN: fmt.Sprintf("%s://key@%s/1", u.Scheme, u.Host),

		ResourceAttributes: []label.KeyValue{
			label.String("resource1", "resource1-value"),
		},
	})

	tracer := otel.Tracer("github.com/your/repo")
	genSpan(ctx, tracer)

	err = upclient.Close()
	require.Nil(t, err)

	require.Equal(t, 1, len(in.Spans))

	s0 := in.Spans[0]

	require.NotZero(t, s0.ID)
	require.Zero(t, s0.ParentID)
	require.NotZero(t, s0.TraceID)
	require.Equal(t, "main span", s0.Name)
	require.Equal(t, "internal", s0.Kind)
	require.NotZero(t, s0.StartTime)
	require.NotZero(t, s0.EndTime)

	set := label.NewSet(s0.Resource...)
	val, ok := set.Value("resource1")
	require.True(t, ok)
	require.Equal(t, "resource1-value", val.AsString())

	require.Equal(t, spanexp.KeyValueSlice{label.String("attr1", "attr1-value")}, s0.Attrs)

	require.Equal(t, "unset", s0.StatusCode)
	require.Equal(t, "", s0.StatusMessage)

	require.Equal(t, 1, len(s0.Events))
	e0 := s0.Events[0]
	require.Equal(t, "event1", e0.Name)
	require.Equal(t, spanexp.KeyValueSlice{label.Int("event1", 123)}, e0.Attrs)

	require.Equal(t, 1, len(s0.Links))
	l0 := s0.Links[0]
	require.NotZero(t, l0.TraceID)
	require.NotZero(t, l0.SpanID)
	require.Equal(t, spanexp.KeyValueSlice{label.Float64("link1", 0.123)}, l0.Attrs)

	require.Equal(t, "github.com/your/repo", s0.Tracer.Name)
	require.Equal(t, "", s0.Tracer.Version)
}

func genSpan(ctx context.Context, tracer trace.Tracer) {
	var traceID [16]byte
	traceID[0] = 0xff

	var spanID [8]byte
	spanID[1] = 0xff

	link1 := trace.Link{
		SpanContext: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		Attributes: []label.KeyValue{label.Float64("link1", 0.123)},
	}

	_, span := tracer.Start(ctx, "main span", trace.WithLinks(link1))

	span.SetAttributes(label.String("attr1", "attr1-value"))
	span.AddEvent("event1", trace.WithAttributes(label.Int("event1", 123)))

	span.End()
}
