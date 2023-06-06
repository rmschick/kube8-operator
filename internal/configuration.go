package internal

import "github.com/FishtechCSOC/locomotive/v8/pkg/types"

type Deployment struct {
	Metadata    types.Metadata `mapstructure:"metadata"`
	Region      string         `mapstructure:"region"`
	Collector   string         `mapstructure:"collector"`
	Instance    string         `mapstructure:"instance"`
	Environment string         `mapstructure:"environment"`
	Namespace   string         `mapstructure:"namespace"`
}
