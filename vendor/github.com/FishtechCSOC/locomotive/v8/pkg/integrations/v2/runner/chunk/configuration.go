package chunk

import (
	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
	"github.com/alecthomas/units"
)

type Configuration struct {
	SizeInBytes int `mapstructure:"sizeInBytes"`
	EntryCount  int `mapstructure:"entryCount"`
}

// nolint: gochecknoglobals, gomnd
var (
	// ConfigurationDefaults defaults for configuration with the prefix of 'chunk'.
	ConfigurationDefaults = configuration.CreateDefaults("chunk").WithFields(map[string]any{
		"sizeInBytes": 10 * units.Megabyte,
	})
)
