package configuration

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configuration interface {
	// Validate should be used to bubble up either ambiguous issues with the given
	// configuration or any other potential
	// misconfiguration.Configuration Instead of returning on first error,
	// a list of errors should be gathered and returned to
	// quicken time to correct.Configuration
	Validate(context.Context) error

	// Defaults is intended for use with spf13/viper which can only load defaults
	// in a specific manor and is only vetted for use with spf13/viper and no
	// other conffiguration library
	Defaults(separator string) map[string]any
}

var _ error = (*ValidationError)(nil)

type ValidationError struct {
	Field  string
	Reason string
}

func NewValidationError(field, reason string) *ValidationError {
	return &ValidationError{
		Field:  field,
		Reason: reason,
	}
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s failed validation due to: %s", v.Field, v.Reason)
}

// GetConfig leverages the viper library to attempt loading in a configuration file.
func GetConfig(viperRef *viper.Viper, configuration any) error {
	initialSetup(viperRef)

	return loadConfig(viperRef, configuration)
}

// GetConfigWithDefaults leverages the viper library to attempt loading in a configuration file while setting defaults.
func GetConfigWithDefaults(viperRef *viper.Viper, configuration any, defaults Defaults) error {
	initialSetup(viperRef)

	for key, value := range defaults.GetFields() {
		viperRef.SetDefault(key, value)
	}

	return loadConfig(viperRef, configuration)
}

func initialSetup(viperRef *viper.Viper) {
	viperRef.SetConfigName("configuration")
	viperRef.SetConfigType("yaml")
	viperRef.AddConfigPath("/etc/integration/")
	viperRef.AddConfigPath("$HOME/integration")
	viperRef.AddConfigPath(".")

	viperRef.AutomaticEnv()
}

func loadConfig(viperRef *viper.Viper, configuration any) error {
	err := viperRef.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read in configuration")
	}

	err = viperRef.Unmarshal(&configuration)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal configuration into expected structure")
	}

	return nil
}
