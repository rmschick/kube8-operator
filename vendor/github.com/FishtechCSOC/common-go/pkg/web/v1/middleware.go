package web

import (
	"github.com/gin-gonic/gin"
)

// The Middleware interface allows custom middleware to be thrown into the server, which are placed after all built-in
// middleware.
type Middleware interface {
	// Name should return some sort of unique ID for the middleware that is also human readable so it can be configured
	// per entrypoint.
	Name() string
	// Handle should wrap whatever is next in the handler chain with whatever business logic it needs.
	Handle(*gin.Context)
}
