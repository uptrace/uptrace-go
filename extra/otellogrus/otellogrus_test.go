package otellogrus

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Test struct {
	log     func(context.Context)
	require func(sdktrace.Event)
}

func TestOtelLogrus(t *testing.T) {
	tests := []Test{
		{
			log: func(ctx context.Context) {
				logrus.WithContext(ctx).Info("hello")
			},
			require: func(event sdktrace.Event) {
				require.Equal(t, []attribute.KeyValue{
					logSeverityKey.String("INFO"),
					logMessageKey.String("hello"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				logrus.WithContext(ctx).WithField("foo", "bar").Warn("hello")
			},
			require: func(event sdktrace.Event) {
				require.Equal(t, []attribute.KeyValue{
					logSeverityKey.String("WARN"),
					logMessageKey.String("hello"),
					attribute.String("foo", "bar"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				err := errors.New("some error")
				logrus.WithContext(ctx).WithError(err).Error("hello")
			},
			require: func(event sdktrace.Event) {
				require.Equal(t, []attribute.KeyValue{
					logSeverityKey.String("ERROR"),
					logMessageKey.String("hello"),
					semconv.ExceptionTypeKey.String("*errors.errorString"),
					semconv.ExceptionMessageKey.String("some error"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				logrus.SetReportCaller(true)
				logrus.WithContext(ctx).Info("hello")
				logrus.SetReportCaller(false)
			},
			require: func(event sdktrace.Event) {
				m := make(map[attribute.Key]attribute.Value, len(event.Attributes))
				for _, kv := range event.Attributes {
					m[kv.Key] = kv.Value
				}

				value, ok := m[semconv.CodeFunctionKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "github.com/uptrace/uptrace-go/extra/otellogrus.TestOtelLogrus")

				value, ok = m[semconv.CodeFilepathKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "otellogrus/otellogrus_test.go")

				_, ok = m[semconv.CodeLineNumberKey]
				require.True(t, ok)
			},
		},
	}

	logrus.AddHook(NewHook(WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))

	for _, test := range tests {
		sr := tracetest.NewSpanRecorder()
		provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
		tracer := provider.Tracer("test")

		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "main")

		test.log(ctx)

		span.End()

		spans := sr.Ended()
		require.Equal(t, 1, len(spans))

		events := spans[0].Events()
		require.Equal(t, 1, len(events))

		event := events[0]
		require.Equal(t, "log", event.Name)
		test.require(event)
	}
}
