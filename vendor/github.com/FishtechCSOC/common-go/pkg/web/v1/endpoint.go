package web

import (
	"github.com/gin-gonic/gin"
)

// The Endpoint interface allows custom endpoints to be thrown into the server, which are placed after all built-in
// endpoints.
type Endpoint interface {
	// Name should return some sort of unique ID for the endpoint that is also human readable so it can be configured
	// per entrypoint.
	Name() string
	// AddRoutes should add the routes it is responsible for to the router group that is passed in.
	AddRoutes(*gin.RouterGroup)
}
