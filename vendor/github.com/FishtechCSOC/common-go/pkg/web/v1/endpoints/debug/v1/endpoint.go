package debug

import (
	httpPProf "net/http/pprof"

	"github.com/gin-gonic/gin"

	"github.com/FishtechCSOC/common-go/pkg/web/v1"
)

const (
	namespace = "debug"
)

var _ web.Endpoint = (*Endpoint)(nil)

// Endpoint is the manager around the endpoint routes.
type Endpoint struct{}

// CreateEndpoint creates a new endpoint instance.
func CreateEndpoint() *Endpoint {
	return &Endpoint{}
}

// Name returns the unique ID for this endpoint.
func (endpoint *Endpoint) Name() string {
	return namespace
}

// AddRoutes adds the meta endpoints to the router group.
func (endpoint *Endpoint) AddRoutes(router *gin.RouterGroup) {
	router.Any("/pprof/cmdline", endpoint.Cmdline)
	router.Any("/pprof/profile", endpoint.Profile)
	router.Any("/pprof/symbol", endpoint.Symbol)
	router.Any("/pprof/trace", endpoint.Trace)
	router.Any("/pprof/allocs", endpoint.Index)
	router.Any("/pprof/block", endpoint.Index)
	router.Any("/pprof/goroutine", endpoint.Index)
	router.Any("/pprof/heap", endpoint.Index)
	router.Any("/pprof/mutex", endpoint.Index)
	router.Any("/pprof/threadcreate", endpoint.Index)
	router.Any("/pprof/", endpoint.Index)
}

// Index is used to expose the built-in pprof HTTP index route.
func (endpoint *Endpoint) Index(ctx *gin.Context) {
	httpPProf.Index(ctx.Writer, ctx.Request)
}

// Cmdline is used to expose the built-in pprof HTTP cmdline route.
func (endpoint *Endpoint) Cmdline(ctx *gin.Context) {
	httpPProf.Cmdline(ctx.Writer, ctx.Request)
}

// Profile is used to expose the built-in pprof HTTP profile route.
func (endpoint *Endpoint) Profile(ctx *gin.Context) {
	httpPProf.Profile(ctx.Writer, ctx.Request)
}

// Symbol is used to expose the built-in pprof HTTP symbol route.
func (endpoint *Endpoint) Symbol(ctx *gin.Context) {
	httpPProf.Symbol(ctx.Writer, ctx.Request)
}

// Trace is used to expose the built-in pprof HTTP trace route.
func (endpoint *Endpoint) Trace(ctx *gin.Context) {
	httpPProf.Trace(ctx.Writer, ctx.Request)
}
