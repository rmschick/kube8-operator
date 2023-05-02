package poll

import (
	"time"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

// BackoffConfiguration stores all the configuration to the poll retry exponential backoff.
type BackoffConfiguration struct {
	Minimum time.Duration `mapstructure:"minimum"`
	Maximum time.Duration `mapstructure:"maximum"`
	Factor  float64       `mapstructure:"factor"`
}

// LifecycleConfiguration stores all configuration related to a poller's lifecycle.
type LifecycleConfiguration struct {
	PollingInterval   time.Duration        `mapstructure:"pollingInterval"`
	InitialWindow     time.Duration        `mapstructure:"initialWindow"`
	WindowOffset      time.Duration        `mapstructure:"windowOffset"`
	MaxWindowSize     time.Duration        `mapstructure:"maxWindowSize"`
	ShutdownTimeout   time.Duration        `mapstructure:"shutdownTimeout"`
	SlidingWindowSize time.Duration        `mapstructure:"slidingWindowSize"`
	Backoff           BackoffConfiguration `mapstructure:"backoff"`
}

// Configuration stores all configuration relevant to an instance of a poller.
type Configuration struct {
	Lifecycle LifecycleConfiguration `mapstructure:"lifecycle"`
	// Deprecated
	RetryCount int `mapstructure:"retryCount"`
}

// nolint: gochecknoglobals, gomnd
var (
	BackoffDefaults = configuration.CreateDefaults("backoff").WithFields(map[string]any{
		"minimum": time.Minute,
		"maximum": time.Minute * 15,
		"factor":  2.0,
	})

	LifecycleDefaults = configuration.CreateDefaults("lifecycle").WithFields(map[string]any{
		"pollingInterval": time.Minute * 1,
		"windowOffset":    time.Minute * 5,
		"initialWindow":   time.Minute * 10,
		"maxWindowSize":   time.Minute * 30,
		"shutdownTimeout": time.Second * 30,
	}).WithChildren(BackoffDefaults)

	// ConfigurationDefaults defaults for configuration with the prefix of 'retriever'.
	ConfigurationDefaults = configuration.CreateDefaults("retriever").WithChildren(LifecycleDefaults)
)
