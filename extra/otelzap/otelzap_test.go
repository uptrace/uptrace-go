package otelzap

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

type Test struct {
	log     func(context.Context, *Logger)
	require func(sdktrace.Event)
}

func TestOtelZap(t *testing.T) {
	tests := []Test{
		{
			log: func(ctx context.Context, log *Logger) {
				log.Ctx(ctx).Info("hello")
			},
			require: func(event sdktrace.Event) {
				require.Equal(t, []attribute.KeyValue{
					logSeverityKey.String("INFO"),
					logMessageKey.String("hello"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context, log *Logger) {
				log.InfoContext(ctx, "hello")
			},
			require: func(event sdktrace.Event) {
				require.Equal(t, []attribute.KeyValue{
					logSeverityKey.String("INFO"),
					logMessageKey.String("hello"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context, log *Logger) {
				log.Ctx(ctx).Warn("hello", zap.String("foo", "bar"))
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
			log: func(ctx context.Context, log *Logger) {
				err := errors.New("some error")
				log.Ctx(ctx).Error("hello", zap.Error(err))
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
			log: func(ctx context.Context, log *Logger) {
				log = log.Clone(WithStackTrace(true))
				log.Ctx(ctx).Info("hello")
			},
			require: func(event sdktrace.Event) {
				var stackTrace attribute.KeyValue

				for _, attr := range event.Attributes {
					if attr.Key == semconv.ExceptionStacktraceKey {
						stackTrace = attr
						break
					}
				}

				require.Equal(t, semconv.ExceptionStacktraceKey, stackTrace.Key)
				require.NotZero(t, stackTrace.Value.AsString())
			},
		},
	}

	logger := New(zap.L())

	for _, test := range tests {
		sr := tracetest.NewSpanRecorder()
		provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
		tracer := provider.Tracer("test")

		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "main")

		test.log(ctx, logger)

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
