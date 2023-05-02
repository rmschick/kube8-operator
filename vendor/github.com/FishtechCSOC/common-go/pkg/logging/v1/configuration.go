package logging

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

type Configuration struct {
	Format       string `mapstructure:"format"`
	Prefix       string `mapstructure:"prefix"`
	Verbose      bool   `mapstructure:"verbose"`
	OmitMetadata bool   `mapstructure:"omitMetadata"`
}

// nolint: gochecknoglobals
var (
	// ConfigurationDefaults defaults for configuration with the prefix of 'logging'.
	ConfigurationDefaults = configuration.CreateDefaults("logging").WithFields(map[string]any{
		"format":       JSONFormat,
		"verbose":      false,
		"omitMetadata": false,
		"prefix":       build.Program,
	})
)

func (c *Configuration) Validate(_ context.Context) error {
	var err error
	if strings.TrimSpace(c.Format) == "" {
		err = multierror.Append(err, configuration.NewValidationError("format", " cannot be blank"))
	}

	if strings.TrimSpace(c.Prefix) == "" {
		err = multierror.Append(err, configuration.NewValidationError("prefix", " cannot be blank"))
	}

	return errors.Wrap(err, "validation failed for configuration")
}
