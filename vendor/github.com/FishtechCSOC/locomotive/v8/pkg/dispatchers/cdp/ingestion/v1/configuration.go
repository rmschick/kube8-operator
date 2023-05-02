package ingestion

import (
	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

// Configuration contains all the necessary configuration for sending processed alerts to logstash.
type Configuration struct{}

// nolint: gochecknoglobals
var (
	// ConfigurationDefaults defaults for configuration with the prefix of 'dispatcher'.
	ConfigurationDefaults = configuration.CreateDefaults("dispatcher")
)
