package stream

import (
	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

type Configuration struct{}

// nolint: gochecknoglobals
var (
	// ConfigurationDefaults defaults for configuration with the prefix of 'retriever'.
	ConfigurationDefaults = configuration.CreateDefaults("retriever")
)
