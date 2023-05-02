package http

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/middleware/accesslog/v1"
)

// Configuration stores all configuration relevant to a single entrypoint.
type Configuration struct {
	EntrypointType    string                  `mapstructure:"type"`
	Port              int                     `mapstructure:"port"`
	Host              string                  `mapstructure:"host"`
	ReadTimeout       time.Duration           `mapstructure:"readTimeout"`
	ReadHeaderTimeout time.Duration           `mapstructure:"readHeaderTimeout"`
	WriteTimeout      time.Duration           `mapstructure:"writeTimeout"`
	IdleTimeout       time.Duration           `mapstructure:"idleTimeout"`
	Endpoints         []string                `mapstructure:"endpoints"`
	Middleware        []string                `mapstructure:"middleware"`
	Metrics           MetricsConfiguration    `mapstructure:"metrics"`
	Accesslog         accesslog.Configuration `mapstructure:"accesslog"`
}

// MetricsConfiguration stores all configuration relevant to the metrics middleware.
type MetricsConfiguration struct {
	Enabled bool `mapstructure:"enabled"`
}

// nolint: gochecknoglobals, gomnd
var (
	ConfigurationDefaults = configuration.CreateDefaults("http").WithFields(map[string]any{
		"type":              "http",
		"port":              8888,
		"host":              "0.0.0.0",
		"readTimeout":       time.Second * 5,
		"readHeaderTimeout": time.Second * 5,
		"writeTimeout":      time.Second * 5,
		"idleTimeout":       time.Second * 5,
		"endpoints":         []string{"meta", "metrics", "debug"},
		"middleware":        []string{},
	})
)

func (c *MetricsConfiguration) Validate(_ context.Context) error {
	return nil
}

func (c *Configuration) Validate(_ context.Context) error {
	var err error

	if strings.TrimSpace(c.EntrypointType) == "" {
		err = multierror.Append(err, configuration.NewValidationError("entry point", " cannot be blank"))
	}

	if c.Port < 0 || c.Port > math.MaxInt16 {
		err = multierror.Append(err, configuration.NewValidationError("port", " must be between 0 and 2^16"))
	}

	if strings.TrimSpace(c.Host) == "" {
		err = multierror.Append(err, configuration.NewValidationError("host", "cannot be blank"))
	}

	return errors.Wrap(err, "validation failed for configuration")
}
