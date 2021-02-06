package otelzap

import (
	"go.uber.org/zap/zapcore"
)

// Option applies a configuration to the given config.
type Option interface {
	Apply(*OtelCore)
}

// optionFunc is a function type that applies a particular
// configuration to the logrus hook.
type optionFunc func(core *OtelCore)

// Apply will apply the option to the logrus hook.
func (o optionFunc) Apply(core *OtelCore) {
	o(core)
}

func WithLevel(enab zapcore.LevelEnabler) Option {
	return optionFunc(func(core *OtelCore) {
		core.LevelEnabler = enab
	})
}

// WithErrorStatusLevel sets the maximum logrus logging level on which
// the span status is set to codes.Error.
//
// The default is <= logrus.ErrorLevel.
func WithErrorStatusLevel(level zapcore.Level) Option {
	return optionFunc(func(core *OtelCore) {
		core.errorStatusLevel = level
	})
}
