package firestore

import (
	"time"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
)

type Configuration struct {
	ProjectID  string        `mapstructure:"projectID"`
	TimeToLive time.Duration `mapstructure:"timeToLive"`
}

// nolint: gochecknoglobals, gomnd
var (
	ConfigurationDefaults = configuration.CreateDefaults("cursor").WithFields(map[string]interface{}{
		"timeToLive": time.Hour * 24 * 30, // 30 days
	})
)
