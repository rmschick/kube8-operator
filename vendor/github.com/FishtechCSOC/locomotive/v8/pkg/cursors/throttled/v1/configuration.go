package throttled

import (
	"context"
	"time"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

var _ configuration.Configuration = (*Configuration)(nil)

type Configuration struct {
	ThrottleOffset time.Duration `mapstructure:"throttleOffset" yaml:"throttleOffset" json:"throttleOffset"`
}

func (c Configuration) Validate(_ context.Context) error {
	return nil
}

func (c Configuration) Defaults(_ string) map[string]any {
	return map[string]any{}
}
