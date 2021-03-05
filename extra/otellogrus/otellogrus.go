package otellogrus

import (
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")

	codeFunctionKey = attribute.Key("code.function")
	codeFilepathKey = attribute.Key("code.filepath")
	codeLinenoKey   = attribute.Key("code.lineno")

	exceptionTypeKey    = attribute.Key("exception.type")
	exceptionMessageKey = attribute.Key("exception.message")
)

// Option applies a configuration to the given config.
type Option interface {
	Apply(*LoggingHook)
}

// optionFunc is a function type that applies a particular
// configuration to the logrus hook.
type optionFunc func(hook *LoggingHook)

// Apply will apply the option to the logrus hook.
func (o optionFunc) Apply(hook *LoggingHook) {
	o(hook)
}

// WithLevels sets the logrus logging levels on which the hook is fired.
//
// The default is all levels between logrus.PanicLevel and logrus.WarnLevel inclusive.
func WithLevels(levels ...logrus.Level) Option {
	return optionFunc(func(hook *LoggingHook) {
		hook.levels = levels
	})
}

// WithErrorStatusLevel sets the maximum logrus logging level on which
// the span status is set to codes.Error.
//
// The default is <= logrus.ErrorLevel.
func WithErrorStatusLevel(level logrus.Level) Option {
	return optionFunc(func(hook *LoggingHook) {
		hook.errorStatusLevel = level
	})
}

// LoggingHook is a logrus hook that adds logs to the active span as events.
type LoggingHook struct {
	levels           []logrus.Level
	errorStatusLevel logrus.Level
}

var _ logrus.Hook = (*LoggingHook)(nil)

// NewLoggingHook returns a logrus hook.
func NewLoggingHook(opts ...Option) *LoggingHook {
	hook := &LoggingHook{
		levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
		errorStatusLevel: logrus.ErrorLevel,
	}

	for _, opt := range opts {
		opt.Apply(hook)
	}

	return hook
}

// Fire is a logrus hook that is fired on a new log entry.
func (hook *LoggingHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(entry.Data)+2+3)

	attrs = append(attrs, logSeverityKey.String(levelString(entry.Level)))
	attrs = append(attrs, logMessageKey.String(entry.Message))

	if entry.Caller != nil {
		if entry.Caller.Function != "" {
			attrs = append(attrs, codeFunctionKey.String(entry.Caller.Function))
		}
		if entry.Caller.File != "" {
			attrs = append(attrs, codeFilepathKey.String(entry.Caller.File))
			attrs = append(attrs, codeLinenoKey.Int(entry.Caller.Line))
		}
	}

	for k, v := range entry.Data {
		if k == "error" {
			if err, ok := v.(error); ok {
				typ := reflect.TypeOf(err).String()
				attrs = append(attrs, exceptionTypeKey.String(typ))
				attrs = append(attrs, exceptionMessageKey.String(err.Error()))
				continue
			}
		}

		attrs = append(attrs, attribute.Any(k, v))
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))

	if entry.Level <= hook.errorStatusLevel {
		span.SetStatus(codes.Error, entry.Message)
	}

	return nil
}

// Levels returns logrus levels on which this hook is fired.
func (hook *LoggingHook) Levels() []logrus.Level {
	return hook.levels
}

func levelString(lvl logrus.Level) string {
	s := lvl.String()
	if s == "warning" {
		s = "warn"
	}
	return strings.ToUpper(s)
}
