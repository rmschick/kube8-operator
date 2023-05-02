package integrations

import (
	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/sirupsen/logrus"
)

// SetupRetrieverLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupRetrieverLogger(retrieverName string, retriever Retriever, logger *logrus.Entry) *logrus.Entry {
	return logging.CreateTypeLogger(logger, retrieverName, retriever)
}

// SetupDispatcherLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupDispatcherLogger(dispatcherName string, dispatcher Dispatcher, logger *logrus.Entry) *logrus.Entry {
	return logging.CreateTypeLogger(logger, dispatcherName, dispatcher)
}

// SetupCursorLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupCursorLogger(cursorName string, cursor Cursor, logger *logrus.Entry) *logrus.Entry {
	return logging.CreateTypeLogger(logger, cursorName, cursor)
}
