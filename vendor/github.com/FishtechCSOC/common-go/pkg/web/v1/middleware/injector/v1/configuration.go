package injector

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	Logger *logrus.Entry
	Router *gin.Engine
}

func (c *Configuration) Validate(_ context.Context) error {
	var err error

	return err
}
