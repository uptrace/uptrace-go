package otelzap

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logSeverityKey = label.Key("log.severity")
	logMessageKey  = label.Key("log.message")

	exceptionTypeKey    = label.Key("exception.type")
	exceptionMessageKey = label.Key("exception.message")
	exceptionStacktrace = label.Key("exception.stacktrace")

	codeFunctionKey = label.Key("code.function")
	codeFilepathKey = label.Key("code.filepath")
	codeLinenoKey   = label.Key("code.lineno")
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

	attrs := make([]label.KeyValue, 0, len(fields)+2+3+1)

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

func appendField(attrs []label.KeyValue, f zapcore.Field) []label.KeyValue {
	switch f.Type {
	case zapcore.BoolType:
		attr := label.Bool(f.Key, f.Integer == 1)
		return append(attrs, attr)

	case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int64Type:
		attr := label.Int64(f.Key, f.Integer)
		return append(attrs, attr)
	case zapcore.Int32Type:
		attr := label.Int32(f.Key, int32(f.Integer))
		return append(attrs, attr)

	case zapcore.Uint32Type:
		attr := label.Uint32(f.Key, uint32(f.Integer))
		return append(attrs, attr)
	case zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint64Type, zapcore.UintptrType:
		attr := label.Uint64(f.Key, uint64(f.Integer))
		return append(attrs, attr)

	case zapcore.Float32Type:
		attr := label.Float32(f.Key, math.Float32frombits(uint32(f.Integer)))
		return append(attrs, attr)
	case zapcore.Float64Type:
		attr := label.Float64(f.Key, math.Float64frombits(uint64(f.Integer)))
		return append(attrs, attr)

	case zapcore.Complex64Type:
		s := strconv.FormatComplex(complex128(f.Interface.(complex64)), 'E', -1, 64)
		attr := label.String(f.Key, s)
		return append(attrs, attr)
	case zapcore.Complex128Type:
		s := strconv.FormatComplex(f.Interface.(complex128), 'E', -1, 128)
		attr := label.String(f.Key, s)
		return append(attrs, attr)

	case zapcore.StringType:
		attr := label.String(f.Key, f.String)
		return append(attrs, attr)
	case zapcore.BinaryType, zapcore.ByteStringType:
		attr := label.String(f.Key, string(f.Interface.([]byte)))
		return append(attrs, attr)
	case zapcore.StringerType:
		attr := label.String(f.Key, f.Interface.(fmt.Stringer).String())
		return append(attrs, attr)

	case zapcore.DurationType, zapcore.TimeType:
		attr := label.Int64(f.Key, f.Integer)
		return append(attrs, attr)
	case zapcore.TimeFullType:
		attr := label.Int64(f.Key, f.Interface.(time.Time).UnixNano())
		return append(attrs, attr)
	case zapcore.ErrorType:
		err := f.Interface.(error)
		typ := reflect.TypeOf(err).String()
		attrs = append(attrs, exceptionTypeKey.String(typ))
		attrs = append(attrs, exceptionMessageKey.String(err.Error()))
		return attrs
	case zapcore.ReflectType:
		attr := label.Any(f.Key, f.Interface)
		return append(attrs, attr)
	case zapcore.SkipType:
		return attrs

	case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType:
		return attrs

	default:
		attr := label.String(f.Key+"_error", fmt.Sprintf("unknown field type: %v", f))
		return append(attrs, attr)
	}
}

func levelString(lvl zapcore.Level) string {
	if lvl == zapcore.DPanicLevel {
		return "PANIC"
	}
	return lvl.CapitalString()
}
