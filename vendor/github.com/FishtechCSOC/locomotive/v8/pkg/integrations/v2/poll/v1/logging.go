package poll

import (
	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/sirupsen/logrus"
)

// SetupPollerLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupPollerLogger(pollerName string, poller Poller, logger *logrus.Entry) *logrus.Entry {
	return logging.CreateTypeLogger(logger, pollerName, poller)
}
