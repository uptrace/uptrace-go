package otelzap

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")

	exceptionTypeKey    = attribute.Key("exception.type")
	exceptionMessageKey = attribute.Key("exception.message")
	exceptionStacktrace = attribute.Key("exception.stacktrace")

	codeFunctionKey = attribute.Key("code.function")
	codeFilepathKey = attribute.Key("code.filepath")
	codeLinenoKey   = attribute.Key("code.lineno")
)

func Wrap(logger *zap.Logger, opts ...Option) *zap.Logger {
	return logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		otelCore := NewOtelCore(opts...)
		return zapcore.NewTee(c, otelCore)
	}))
}

type OtelCore struct {
	zapcore.LevelEnabler
	errorStatusLevel zapcore.Level
}

var _ zapcore.Core = (*OtelCore)(nil)

func NewOtelCore(opts ...Option) *OtelCore {
	core := &OtelCore{
		LevelEnabler:     zap.NewAtomicLevelAt(zap.ErrorLevel),
		errorStatusLevel: zapcore.ErrorLevel,
	}
	for _, opt := range opts {
		opt.Apply(core)
	}
	return core
}

func (c *OtelCore) With(fields []zapcore.Field) zapcore.Core {
	return c
}

func (c *OtelCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *OtelCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if ent.Ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ent.Ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(fields)+2+3+1)

	attrs = append(attrs, logSeverityKey.String(levelString(ent.Level)))
	attrs = append(attrs, logMessageKey.String(ent.Message))

	if ent.Caller.Defined {
		if ent.Caller.Function != "" {
			attrs = append(attrs, codeFunctionKey.String(ent.Caller.Function))
		}
		if ent.Caller.File != "" {
			attrs = append(attrs, codeFilepathKey.String(ent.Caller.File))
			attrs = append(attrs, codeLinenoKey.Int(ent.Caller.Line))
		}
	}

	if ent.Stack != "" {
		attrs = append(attrs, exceptionStacktrace.String(ent.Stack))
	}

	for _, f := range fields {
		if f.Type == zapcore.NamespaceType {
			// should this be a prefix?
			continue
		}
		attrs = appendField(attrs, f)
	}

	span.AddEvent("log",
		trace.WithTimestamp(ent.Time),
		trace.WithAttributes(attrs...))

	if ent.Level >= c.errorStatusLevel {
		span.SetStatus(codes.Error, ent.Message)
	}

	return nil
}

func (c *OtelCore) Sync() error {
	return nil
}

func appendField(attrs []attribute.KeyValue, f zapcore.Field) []attribute.KeyValue {
	switch f.Type {
	case zapcore.BoolType:
		attr := attribute.Bool(f.Key, f.Integer == 1)
		return append(attrs, attr)

	case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type,
		zapcore.Uint32Type, zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint64Type,
		zapcore.UintptrType:
		attr := attribute.Int64(f.Key, f.Integer)
		return append(attrs, attr)

	case zapcore.Float32Type, zapcore.Float64Type:
		attr := attribute.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
		return append(attrs, attr)

	case zapcore.Complex64Type:
		s := strconv.FormatComplex(complex128(f.Interface.(complex64)), 'E', -1, 64)
		attr := attribute.String(f.Key, s)
		return append(attrs, attr)
	case zapcore.Complex128Type:
		s := strconv.FormatComplex(f.Interface.(complex128), 'E', -1, 128)
		attr := attribute.String(f.Key, s)
		return append(attrs, attr)

	case zapcore.StringType:
		attr := attribute.String(f.Key, f.String)
		return append(attrs, attr)
	case zapcore.BinaryType, zapcore.ByteStringType:
		attr := attribute.String(f.Key, string(f.Interface.([]byte)))
		return append(attrs, attr)
	case zapcore.StringerType:
		attr := attribute.String(f.Key, f.Interface.(fmt.Stringer).String())
		return append(attrs, attr)

	case zapcore.DurationType, zapcore.TimeType:
		attr := attribute.Int64(f.Key, f.Integer)
		return append(attrs, attr)
	case zapcore.TimeFullType:
		attr := attribute.Int64(f.Key, f.Interface.(time.Time).UnixNano())
		return append(attrs, attr)
	case zapcore.ErrorType:
		err := f.Interface.(error)
		typ := reflect.TypeOf(err).String()
		attrs = append(attrs, exceptionTypeKey.String(typ))
		attrs = append(attrs, exceptionMessageKey.String(err.Error()))
		return attrs
	case zapcore.ReflectType:
		attr := attribute.Any(f.Key, f.Interface)
		return append(attrs, attr)
	case zapcore.SkipType:
		return attrs

	case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType:
		return attrs

	default:
		attr := attribute.String(f.Key+"_error", fmt.Sprintf("unknown field type: %v", f))
		return append(attrs, attr)
	}
}

func levelString(lvl zapcore.Level) string {
	if lvl == zapcore.DPanicLevel {
		return "PANIC"
	}
	return lvl.CapitalString()
}
