package zapdefaults

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

// Option is a type that can be used to override the logger configuration.
type Option interface {
	apply(*zap.Config) error
}

// OptionFunc is a function that conforms to the Option interface.
type OptionFunc func(config *zap.Config) error

// apply conforms to the Option interface and applies the OptionFunc.
func (f OptionFunc) apply(config *zap.Config) error {
	return f(config)
}

// EnvironmentOverrides supports configuring zap with environment variables.
// The environment is parsed to overwrite the configuration struct using
// github.com/kelseyhightower/envconfig with a default prefix of "ZAP".
//
// For example, to change the log level, you may set:
//     ZAP_LEVEL=warn
//
// Or to change the encoding,
//     ZAP_ENCODING=json
//
// Note: selecting a development or production as the mode does more than
// just choosing the encoding, so ensure that if you override all of the
// relevant settings from the default configuration if you want to control
// this manually.
func EnvironmentOverrides(prefix string) Option {
	if prefix == "" {
		prefix = "ZAP"
	}

	return OptionFunc(func(cfg *zap.Config) error {
		if err := envconfig.Process(prefix, cfg); err != nil {
			return err
		}

		return nil
	})
}
