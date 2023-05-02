package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
)

const (
	WebFieldKey        = "web"
	EndpointFieldKey   = WebFieldKey + ".endpoint"
	EntrypointFieldKey = WebFieldKey + ".entrypoint"
	MiddlewareFieldKey = WebFieldKey + ".middleware"
	ServerFieldKey     = WebFieldKey + ".server"
)

// SetupEndpointLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupEndpointLogger(endpointName string, endpoint any, ctx *gin.Context) *logrus.Entry {
	return logging.CreateTypeLogger(loggerFromContext(ctx), endpointName, endpoint).WithField(EndpointFieldKey, endpointName)
}

// SetupEntrypointLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupEntrypointLogger(logger *logrus.Entry, entrypointName string, entrypoint any) *logrus.Entry {
	return logging.CreateTypeLogger(logger, entrypointName, entrypoint).WithField(EntrypointFieldKey, entrypointName)
}

// SetupMiddlewareLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupMiddlewareLogger(middlewareName string, middleware any, ctx *gin.Context) *logrus.Entry {
	return logging.CreateTypeLogger(loggerFromContext(ctx), middlewareName, middleware).WithField(MiddlewareFieldKey, middleware)
}

// SetupServerLogger is just used to help remove some boilerplate around setting up loggers to make code more DRY.
func SetupServerLogger(logger *logrus.Entry, serverName string, server any) *logrus.Entry {
	return logging.CreateTypeLogger(logger, serverName, server).WithField(ServerFieldKey, serverName)
}

func loggerFromContext(ctx *gin.Context) *logrus.Entry {
	loggerItem, ok := ctx.Get("logger")

	if !ok {
		return nil
	}

	logger, ok := loggerItem.(*logrus.Entry)

	if !ok {
		return nil
	}

	return logger
}
