package zapdefaults

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Preset is an option that specifies a default configuration.
type Preset int32

const (
	Invalid Preset = iota
	// Development uses a console encoder, colored output, and more friendly string-encoded durations.
	Development
	// Development uses a json encoder with durations in seconds.
	Production
	// Dynamic uses a Development preset if attached to a TTY and a Production preset if not.
	Dynamic
)

// String returns a string describing the level.
func (m Preset) String() string {
	switch m {
	case Development:
		return "development"
	case Production:
		return "production"
	case Dynamic:
		return "dynamic"
	default:
		return "invalid"
	}
}

// UnmarshalText initializes a preset from text.
//
// Preset conforms to the encoding.TextUnmarshaler interface,
// which also lets it support unmarshaling values from json
// structures.
func (m *Preset) UnmarshalText(text []byte) error {
	lowerModeString := strings.ToLower(string(text))
	switch lowerModeString {
	case Development.String():
		*m = Development
	case Production.String():
		*m = Production
	case Dynamic.String():
		*m = Dynamic
	default:
		*m = Invalid
		return fmt.Errorf("invalid logger mode: %s", text)
	}
	return nil
}

// MarshalText supports encoding the string as text.
//
// Preset conforms to the encoding.TextMarshaler interface.
func (m Preset) MarshalText() ([]byte, error) {
	s := m.String()
	b := []byte(s)
	if s == "invalid" {
		return b, fmt.Errorf("cannot marshal invalid preset")
	}
	return b, nil
}

func (m Preset) apply(config *zap.Config) error {
	// apply the preset overrides
	switch m {
	case Development:
		// override the current settings to start with a clean slate
		*config = *DefaultConfiguration()
		config.Development = true
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	case Production:
		// override the current settings to start with a clean slate
		*config = *DefaultConfiguration()
		config.Development = false
		config.Encoding = "json"
		config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		config.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	case Dynamic:
		var err error
		if isatty.IsTerminal(os.Stdout.Fd()) {
			err = Development.apply(config)
		} else {
			err = Production.apply(config)
		}
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("cannot apply invalid preset")
	}
	return nil
}
