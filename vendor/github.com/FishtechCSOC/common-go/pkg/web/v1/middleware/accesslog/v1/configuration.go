package accesslog

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

const (
	// Keep is the keep string value.
	Keep = "keep"
	// Drop is the drop string value.
	Drop = "drop"
	// Redact is the redact string value.
	Redact = "redact"
)

type DirectionalConfiguration struct {
	Write       bool                `mapstructure:"write"`
	IncludeBody bool                `mapstructure:"includeBody"`
	Headers     FilterConfiguration `mapstructure:"headers"`
	QueryParams FilterConfiguration `mapstructure:"queryParams"`
}

type FilterConfiguration struct {
	Default  string            `mapstructure:"default"`
	Override map[string]string `mapstructure:"override"`
}

type Configuration struct {
	SamplingRatio float64                  `mapstructure:"samplingRatio"`
	Requests      DirectionalConfiguration `mapstructure:"requests"`
	Responses     DirectionalConfiguration `mapstructure:"responses"`
}

// nolint: gochecknoglobals
var (
	// ConfigurationDefaults defaults for configuration with the prefix of 'accesslog'.
	ConfigurationDefaults = configuration.CreateDefaults("accesslog")
)

func (c *DirectionalConfiguration) Validate(_ context.Context) error {
	var err error

	return err
}

func (c *FilterConfiguration) Validate(_ context.Context) error {
	var err error
	if strings.TrimSpace(c.Default) == "" {
		err = multierror.Append(err, configuration.NewValidationError("default", " cannot be blank"))
	}

	return errors.Wrap(err, "validation failed for configuration")
}

func (c *Configuration) Validate(_ context.Context) error {
	return nil
}
