package otellogrus

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/oteltest"
)

type Test struct {
	log     func(context.Context)
	require func(oteltest.Event)
}

func TestLogrusHook(t *testing.T) {
	tests := []Test{
		{
			log: func(ctx context.Context) {
				logrus.WithContext(ctx).Info("hello")
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[label.Key]label.Value{
					logSeverityKey: label.StringValue("INFO"),
					logMessageKey:  label.StringValue("hello"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				logrus.WithContext(ctx).WithField("foo", "bar").Warn("hello")
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[label.Key]label.Value{
					logSeverityKey:   label.StringValue("WARN"),
					logMessageKey:    label.StringValue("hello"),
					label.Key("foo"): label.StringValue("bar"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				err := errors.New("some error")
				logrus.WithContext(ctx).WithError(err).Error("hello")
			},
			require: func(event oteltest.Event) {
				require.Equal(t, map[label.Key]label.Value{
					logSeverityKey:      label.StringValue("ERROR"),
					logMessageKey:       label.StringValue("hello"),
					exceptionTypeKey:    label.StringValue("*errors.errorString"),
					exceptionMessageKey: label.StringValue("some error"),
				}, event.Attributes)
			},
		},
		{
			log: func(ctx context.Context) {
				logrus.SetReportCaller(true)
				logrus.WithContext(ctx).Info("hello")
				logrus.SetReportCaller(false)
			},
			require: func(event oteltest.Event) {
				value, ok := event.Attributes[codeFunctionKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "github.com/uptrace/uptrace-go/extra/otellogrus")

				value, ok = event.Attributes[codeFilepathKey]
				require.True(t, ok)
				require.Contains(t, value.AsString(), "uptrace-go/extra/otellogrus/otellogrus_test.go")

				_, ok = event.Attributes[codeLinenoKey]
				require.True(t, ok)
			},
		},
	}

	sr := new(oteltest.StandardSpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))

	logrus.AddHook(NewLoggingHook(WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))
	tracer := provider.Tracer("test")

	for _, test := range tests {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "main")

		test.log(ctx)

		events := span.(*oteltest.Span).Events()
		require.Equal(t, 1, len(events))

		event := events[0]
		require.Equal(t, "log", event.Name)
		test.require(event)

		span.End()
	}
}

func TestSpanStatus(t *testing.T) {
	sr := new(oteltest.StandardSpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))

	logrus.AddHook(NewLoggingHook())
	tracer := provider.Tracer("test")

	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "main")

	logrus.WithContext(ctx).Info("hello")
	require.Equal(t, codes.Unset, span.(*oteltest.Span).StatusCode())

	logrus.WithContext(ctx).Error("hello")
	require.Equal(t, codes.Error, span.(*oteltest.Span).StatusCode())
}
