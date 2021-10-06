package otelzap

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")
)

// Logger is a thin wrapper for zap.Logger that adds Ctx method.
type Logger struct {
	*zap.Logger

	stackTrace       bool
	minLevel         zapcore.Level
	errorStatusLevel zapcore.Level
}

// Deprecated. Use New instead.
func Wrap(logger *zap.Logger, opts ...Option) *Logger {
	return New(logger, opts...)
}

func New(logger *zap.Logger, opts ...Option) *Logger {
	l := &Logger{
		Logger: logger,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// WithOptions clones the current Logger, applies the supplied Options,
// and returns the resulting Logger. It's safe to use concurrently.
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	clone := *l
	clone.Logger = l.Logger.WithOptions(opts...)
	return &clone
}

// Clone clones the current logger applying the supplied options.
func (l *Logger) Clone(opts ...Option) *Logger {
	clone := *l
	for _, opt := range opts {
		opt(&clone)
	}
	return &clone
}

// Ctx returns a new logger with context.
func (l *Logger) Ctx(ctx context.Context) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: ctx,
		l:   l,
	}
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.DebugLevel, msg, fields)
	l.Debug(msg, fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.InfoLevel, msg, fields)
	l.Info(msg, fields...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.WarnLevel, msg, fields)
	l.Warn(msg, fields...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.ErrorLevel, msg, fields)
	l.Error(msg, fields...)
}

func (l *Logger) DPanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.DPanicLevel, msg, fields)
	l.DPanic(msg, fields...)
}

func (l *Logger) PanicContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.PanicLevel, msg, fields)
	l.Panic(msg, fields...)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, fields ...zapcore.Field) {
	l.log(ctx, zap.FatalLevel, msg, fields)
	l.Fatal(msg, fields...)
}

func (l *Logger) log(
	ctx context.Context, lvl zapcore.Level, msg string, fields []zapcore.Field,
) {
	if lvl < l.minLevel {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	attrs := make([]attribute.KeyValue, 0, 3+len(fields))

	attrs = append(attrs, logSeverityKey.String(levelString(lvl)))
	attrs = append(attrs, logMessageKey.String(msg))

	if l.stackTrace {
		stackTrace := make([]byte, 2048)
		n := runtime.Stack(stackTrace, false)
		attrs = append(attrs, semconv.ExceptionStacktraceKey.String(string(stackTrace[0:n])))
	}

	for _, f := range fields {
		if f.Type == zapcore.NamespaceType {
			// should this be a prefix?
			continue
		}
		attrs = appendField(attrs, f)
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))

	if lvl >= l.errorStatusLevel {
		span.SetStatus(codes.Error, msg)
	}
}

//------------------------------------------------------------------------------

// LoggerWithCtx is a wrapper for Logger that also carries a context.Context.
type LoggerWithCtx struct {
	ctx context.Context
	l   *Logger
}

func (l LoggerWithCtx) Context() context.Context {
	return l.ctx
}

func (l LoggerWithCtx) Logger() *Logger {
	return l.l
}

func (l LoggerWithCtx) ZapLogger() *zap.Logger {
	return l.l.Logger
}

// WithOptions clones the current Logger, applies the supplied Options,
// and returns the resulting Logger. It's safe to use concurrently.
func (l LoggerWithCtx) WithOptions(opts ...zap.Option) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.WithOptions(opts...),
	}
}

// Clone clones the current logger applying the supplied options.
func (l LoggerWithCtx) Clone(opts ...Option) LoggerWithCtx {
	return LoggerWithCtx{
		ctx: l.ctx,
		l:   l.l.Clone(opts...),
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Debug(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.DebugLevel, msg, fields)
	l.l.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Info(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.InfoLevel, msg, fields)
	l.l.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Warn(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.WarnLevel, msg, fields)
	l.l.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l LoggerWithCtx) Error(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.ErrorLevel, msg, fields)
	l.l.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l LoggerWithCtx) DPanic(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.DPanicLevel, msg, fields)
	l.l.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l LoggerWithCtx) Panic(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.PanicLevel, msg, fields)
	l.l.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l LoggerWithCtx) Fatal(msg string, fields ...zapcore.Field) {
	l.l.log(l.ctx, zap.FatalLevel, msg, fields)
	l.l.Fatal(msg, fields...)
}

//------------------------------------------------------------------------------

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
		attrs = append(attrs, semconv.ExceptionTypeKey.String(typ))
		attrs = append(attrs, semconv.ExceptionMessageKey.String(err.Error()))
		return attrs
	case zapcore.ReflectType:
		attr := attrAny(f.Key, f.Interface)
		return append(attrs, attr)
	case zapcore.SkipType:
		return attrs

	case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType:
		return attrs

	default:
		attr := attribute.String(f.Key+"_error", fmt.Sprintf("otelzap: unknown field type: %v", f))
		return append(attrs, attr)
	}
}

func levelString(lvl zapcore.Level) string {
	if lvl == zapcore.DPanicLevel {
		return "PANIC"
	}
	return lvl.CapitalString()
}

func attrAny(k string, value interface{}) attribute.KeyValue {
	if value == nil {
		return attribute.String(k, "<nil>")
	}

	if stringer, ok := value.(fmt.Stringer); ok {
		return attribute.String(k, stringer.String())
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Array:
		rv = rv.Slice(0, rv.Len())
		fallthrough
	case reflect.Slice:
		switch reflect.TypeOf(value).Elem().Kind() {
		case reflect.Bool:
			return attribute.BoolSlice(k, rv.Interface().([]bool))
		case reflect.Int:
			return attribute.IntSlice(k, rv.Interface().([]int))
		case reflect.Int64:
			return attribute.Int64Slice(k, rv.Interface().([]int64))
		case reflect.Float64:
			return attribute.Float64Slice(k, rv.Interface().([]float64))
		case reflect.String:
			return attribute.StringSlice(k, rv.Interface().([]string))
		default:
			return attribute.KeyValue{Key: attribute.Key(k)}
		}
	case reflect.Bool:
		return attribute.Bool(k, rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return attribute.Int64(k, rv.Int())
	case reflect.Float64:
		return attribute.Float64(k, rv.Float())
	case reflect.String:
		return attribute.String(k, rv.String())
	}
	if b, err := json.Marshal(value); b != nil && err == nil {
		return attribute.String(k, string(b))
	}
	return attribute.String(k, fmt.Sprint(value))
}
