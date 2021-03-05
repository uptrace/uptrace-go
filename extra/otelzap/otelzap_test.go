package otelzap

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/oteltest"
	"go.uber.org/zap"
)

type Test struct {
	log     func(zap.LoggerWithCtx)
	require func(oteltest.Event)
}

func TestOtelCore(t *testing.T) {
	tests := []Test{
		{
			log: func(log zap.LoggerWithCtx) {
				log.Info("hello")
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[attribute.Key]attribute.Value{
					logSeverityKey: attribute.StringValue("INFO"),
					logMessageKey:  attribute.StringValue("hello"),
				}, event.Attributes)
			},
		},
		{
			log: func(log zap.LoggerWithCtx) {
				log.Warn("hello", zap.String("foo", "bar"))
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[attribute.Key]attribute.Value{
					logSeverityKey:       attribute.StringValue("WARN"),
					logMessageKey:        attribute.StringValue("hello"),
					attribute.Key("foo"): attribute.StringValue("bar"),
				}, event.Attributes)
			},
		},
		{
			log: func(log zap.LoggerWithCtx) {
				err := errors.New("some error")
				log.Error("hello", zap.Error(err))
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[attribute.Key]attribute.Value{
					logSeverityKey:      attribute.StringValue("ERROR"),
					logMessageKey:       attribute.StringValue("hello"),
					exceptionTypeKey:    attribute.StringValue("*errors.errorString"),
					exceptionMessageKey: attribute.StringValue("some error"),
				}, event.Attributes)
			},
		},
		{
			log: func(log zap.LoggerWithCtx) {
				log = log.WithOptions(zap.AddCaller())
				log.Info("hello")
			},
			require: func(event oteltest.Event) {
				value, ok := event.Attributes[codeFunctionKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "github.com/uptrace/uptrace-go/extra/otelzap")

				value, ok = event.Attributes[codeFilepathKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "uptrace-go/extra/otelzap/otelzap_test.go")

				_, ok = event.Attributes[codeLinenoKey]
				require.True(t, ok)
			},
		},
	}

	sr := new(oteltest.StandardSpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	tracer := provider.Tracer("test")

	core := NewOtelCore(WithLevel(zap.NewAtomicLevelAt(zap.InfoLevel)))
	log := zap.New(core)

	for _, test := range tests {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "main")

		test.log(log.Ctx(ctx))

		events := span.(*oteltest.Span).Events()
		require.Equal(t, 1, len(events))

		event := events[0]
		require.Equal(t, "log", event.Name)
		test.require(event)

		span.End()
	}
}
