package zapdefaults

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger initializes a new logger.
//
// By default, the development logger will be used if running in a TTY interface,
// and the production logger will be used if not. You may change the mode by passing
// a mode Option of either Development or Production which will apply the default
// settings for that mode.
//
// If no options are provided, EnvironmentOverrides option is also enabled by default
// with a default prefix of "ZAP". See the option documentation for additional details.
func NewLogger(opts ...Option) (*zap.Logger, error) {
	cfg, err := BuildConfiguration(opts...)
	if err != nil {
		return nil, err
	}

	return cfg.Build()
}

// BuildConfiguration supports customizing the way that zap is configured.
//
// To build a configuration struct, BuildConfiguration begins with the
// packages' DefaultConfiguration() and then applies each option in order
// to modify the default config. If any option returns an error, it
// is returned directly.
//
// If a Preset (Development or Production) is set, it must be passed
// as the first option, otherwise an error is returned.
//
// See the NewLogger documentation for a discussion of the default
// behaviors when no options are provided.
func BuildConfiguration(opts ...Option) (*zap.Config, error) {
	if len(opts) == 0 {
		opts = append(opts, Dynamic, EnvironmentOverrides(""))
	}

	cfg := DefaultConfiguration()
	for i, opt := range opts {
		if _, isMode := opt.(Preset); isMode && i > 0 {
			return nil, errors.Errorf("logger mode option can only be applied as the first option")
		}

		if err := opt.apply(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// DefaultConfiguration returns a copy of the default logger settings.
func DefaultConfiguration() *zap.Config {
	return &zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{},
	}
}

// MustLogger returns a logger or panics if it fails to be constructed.
func MustLogger(opts ...Option) *zap.Logger {
	logger, err := NewLogger(opts...)
	if err != nil {
		panic(err)
	}
	return logger
}
