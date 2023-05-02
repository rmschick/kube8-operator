package server

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/entrypoints/http/v1"
)

// LifecycleConfiguration stores all configuration related to a server's lifecycle.
type LifecycleConfiguration struct {
	DrainTimeout    time.Duration `mapstructure:"drainTimeout"`
	GracefulTimeout time.Duration `mapstructure:"gracefulTimeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

// EntrypointConfiguration stores all configuration relevant to a single entrypoint.
type EntrypointConfiguration struct {
	HTTP *http.Configuration `mapstructure:"http"`
}

// MetricsConfiguration stores all configuration relevant to the metrics middleware.
type MetricsConfiguration struct {
	Enabled bool `mapstructure:"enabled"`
}

// Configuration stores all configuration relevant to an instance of a server.
type Configuration struct {
	Lifecycle   LifecycleConfiguration             `mapstructure:"lifecycle"`
	Entrypoints map[string]EntrypointConfiguration `mapstructure:"entrypoints"`
}

// nolint: gochecknoglobals, gomnd
var (
	LifecycleConfigurationDefaults = configuration.CreateDefaults("lifecycle").WithFields(map[string]any{
		"drainTimeout":    time.Second * 15,
		"gracefulTimeout": time.Second * 15,
		"shutdownTimeout": time.Second * 5,
	})
	EntrypointsConfiguration = configuration.CreateDefaults("entrypoints").WithChildren(
		configuration.CreateDefaults("introspection").WithChildren(http.ConfigurationDefaults),
	)
	// ConfigurationDefaults defaults for configuration with the prefix of 'server'.
	ConfigurationDefaults = configuration.CreateDefaults("server").WithChildren(
		LifecycleConfigurationDefaults,
		EntrypointsConfiguration,
	)
)

func (c *LifecycleConfiguration) Validate(_ context.Context) error {
	var err error

	return err
}

func (c *EntrypointConfiguration) Validate(_ context.Context) error {
	var err error

	switch {
	case c.HTTP != nil:
	default:
		err = multierror.Append(err, errors.New("entrypoint configuration empty"))
	}

	return errors.Wrap(err, "validation failed for configuration")
}

func (c *MetricsConfiguration) Validate(_ context.Context) error {
	return nil
}

func (c *Configuration) Validate(_ context.Context) error {
	return nil
}
