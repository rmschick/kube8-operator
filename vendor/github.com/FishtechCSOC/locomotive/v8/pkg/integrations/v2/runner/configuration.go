package runner

import (
	"runtime"
	"time"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

type Batch struct {
	SizeInBytes int           `mapstructure:"sizeInBytes"`
	EntryCount  int           `mapstructure:"entryCount"`
	WaitTime    time.Duration `mapstructure:"waitTime"`
}

type Configuration struct {
	WorkerCount  int           `mapstructure:"workerCount"`
	DrainTimeout time.Duration `mapstructure:"drainTimeout"`
	Batch        Batch         `mapstructure:"batch"`
}

// nolint: gochecknoglobals, gomnd
var (
	BatchDefaults = configuration.CreateDefaults("batch")

	// ConfigurationDefaults defaults for configuration with the prefix of 'runner'.
	ConfigurationDefaults = configuration.CreateDefaults("runner").WithFields(map[string]any{
		"workerCount":  runtime.NumCPU(),
		"drainTimeout": time.Second * 15,
	}).WithChildren(BatchDefaults)
)
