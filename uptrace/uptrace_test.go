package uptrace_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/uptrace"
)

func TestInvalidDSN(t *testing.T) {
	t.Skip("overwrites global tracer provider")

	var logger Logger
	uptrace.SetLogger(&logger)

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		DSN: "dsn",
	})

	require.Equal(t,
		`Uptrace is disabled: DSN="dsn" does not have a token`,
		logger.Message())
}

func TestUnknownToken(t *testing.T) {
	t.Skip("overwrites global tracer provider")

	ctx := context.Background()

	var logger Logger
	uptrace.SetLogger(&logger)

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		DSN: "https://UNKNOWN@api.uptrace.dev/2",
	})

	uptrace.ReportError(ctx, errors.New("hello"))
	err := uptrace.Shutdown(ctx)
	require.NoError(t, err)

	require.Equal(t,
		`send failed: status=403: project with such id and token not found (DSN="https://UNKNOWN@api.uptrace.dev/2")`,
		logger.Message())
}

func TestBeforeSpanSend(t *testing.T) {
	ctx := context.Background()

	var got *spanexp.Span

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		DSN: "https://token@api.uptrace.dev/1",

		ServiceName:    "test-filters",
		ServiceVersion: "1.0.0",

		BeforeSpanSend: func(span *spanexp.Span) {
			got = span
		},
	})

	tracer := otel.Tracer("github.com/your/repo")
	_, span := tracer.Start(ctx, "main span")
	span.End()

	err := uptrace.Shutdown(ctx)
	require.NoError(t, err)

	require.NotNil(t, got)
	require.Equal(t, "main span", got.Name)

	set := attribute.NewSet(got.Resource...)
	val, ok := set.Value(semconv.ServiceNameKey)
	require.True(t, ok)
	require.Equal(t, "test-filters", val.AsString())

	val, ok = set.Value(semconv.ServiceVersionKey)
	require.True(t, ok)
	require.Equal(t, "1.0.0", val.AsString())
}

func TestExporter(t *testing.T) {
	ctx := context.Background()

	type In struct {
		Spans   []spanexp.Span `msgpack:"spans"`
		Sampler string         `msgpack:"sampler"`
	}

	var in In

	handler := func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, "application/msgpack", req.Header.Get("Content-Type"))
		require.Equal(t, "zstd", req.Header.Get("Content-Encoding"))

		zr, err := zstd.NewReader(req.Body)
		require.NoError(t, err)
		defer zr.Close()

		dec := msgpack.NewDecoder(zr)
		err = dec.Decode(&in)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		DSN: fmt.Sprintf("%s://key@%s/1", u.Scheme, u.Host),

		ResourceAttributes: []attribute.KeyValue{
			attribute.String("resource1", "resource1-value"),
		},

		Sampler: sdktrace.AlwaysSample(),
	})

	tracer := otel.Tracer("github.com/your/repo")
	genSpan(ctx, tracer)
	err = uptrace.Shutdown(ctx)
	require.NoError(t, err)

	require.Equal(t, "AlwaysOnSampler", in.Sampler)
	require.Equal(t, 1, len(in.Spans))

	s0 := in.Spans[0]

	require.NotZero(t, s0.ID)
	require.Zero(t, s0.ParentID)
	require.NotZero(t, s0.TraceID)
	require.Equal(t, "main span", s0.Name)
	require.Equal(t, "server", s0.Kind)
	require.NotZero(t, s0.StartTime)
	require.NotZero(t, s0.EndTime)

	set := attribute.NewSet(s0.Resource...)

	for _, attrKey := range []string{
		"host.name",
		"telemetry.sdk.language",
		"telemetry.sdk.version",
	} {
		_, ok := set.Value(attribute.Key(attrKey))
		require.True(t, ok)
	}

	val, ok := set.Value("resource1")
	require.True(t, ok)
	require.Equal(t, "resource1-value", val.AsString())

	require.Equal(t, spanexp.KeyValueSlice{attribute.String("attr1", "attr1-value")}, s0.Attrs)

	require.Equal(t, "unset", s0.StatusCode)
	require.Equal(t, "", s0.StatusMessage)

	require.Equal(t, 1, len(s0.Events))
	e0 := s0.Events[0]
	require.Equal(t, "event1", e0.Name)
	require.Equal(t, spanexp.KeyValueSlice{attribute.Int("event1", 123)}, e0.Attrs)

	require.Equal(t, 1, len(s0.Links))
	l0 := s0.Links[0]
	require.NotZero(t, l0.TraceID)
	require.NotZero(t, l0.SpanID)
	require.Equal(t, spanexp.KeyValueSlice{attribute.Float64("link1", 0.123)}, l0.Attrs)

	require.Equal(t, "github.com/your/repo", s0.TracerName)
	require.Equal(t, "", s0.TracerVersion)
}

func genSpan(ctx context.Context, tracer trace.Tracer) {
	var traceID [16]byte
	traceID[0] = 0xff

	var spanID [8]byte
	spanID[1] = 0xff

	link1 := trace.Link{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}),
		Attributes: []attribute.KeyValue{attribute.Float64("link1", 0.123)},
	}

	_, span := tracer.Start(ctx, "main span",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithLinks(link1))

	span.SetAttributes(attribute.String("attr1", "attr1-value"))
	span.AddEvent("event1", trace.WithAttributes(attribute.Int("event1", 123)))

	span.End()
}

//------------------------------------------------------------------------------

type Logger struct {
	msgs []string
}

func (l *Logger) Printf(ctx context.Context, format string, args ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
}

func (l *Logger) Message() string {
	if len(l.msgs) == 0 {
		return ""
	}
	return l.msgs[len(l.msgs)-1]
}
