package injector

import (
	"github.com/gin-gonic/gin"
)

const (
	namespace = "injector"
)

type Middleware struct {
	configuration Configuration
}

func CreateMiddleware(configuration Configuration) *Middleware {
	return &Middleware{
		configuration: configuration,
	}
}

func (middleware *Middleware) Name() string {
	return namespace
}

func (middleware *Middleware) Handle(ctx *gin.Context) {
	// Should find a better way than passing the objects themselves to handle this info
	// Essentially we want to pass metadata about the entrypoint/server for various usage
	ctx.Set("logger", middleware.configuration.Logger)
	ctx.Set("router", middleware.configuration.Router)
}
