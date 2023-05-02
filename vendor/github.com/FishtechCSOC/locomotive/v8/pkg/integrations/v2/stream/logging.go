package stream

import (
	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
)

// SetupStreamerLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupStreamerLogger(streamerName string, streamer integrations.Retriever, logger *logrus.Entry) *logrus.Entry {
	return logging.CreateTypeLogger(logger, streamerName, streamer)
}
