package uptrace_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/uptrace-go/spanexp"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
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

	require.Equal(t, "main span", got.Name)

	set := label.NewSet(got.Resource...)
	val, ok := set.Value(semconv.ServiceNameKey)
	require.True(t, ok)
	require.Equal(t, "test-filters", val.AsString())
}
