package setup

import (
	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
	"github.com/spf13/viper"
)

// BuildConfiguration builds a Viper config with defaults and env vars.
func BuildConfiguration(viperRef *viper.Viper, config any, defaults configuration.Defaults, envVars ...string) {
	for _, envVar := range envVars {
		_ = viperRef.BindEnv(envVar)
	}

	err := configuration.GetConfigWithDefaults(viperRef, config, defaults)
	if err != nil {
		panic(err)
	}
}
