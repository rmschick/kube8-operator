package meta

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tevino/abool"

	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/FishtechCSOC/common-go/pkg/web/v1"
)

const (
	namespace = "meta"
)

var _ web.Endpoint = (*Endpoint)(nil)

// Endpoint is the manager around the endpoint routes.
type Endpoint struct {
	readyChecks []*abool.AtomicBool
}

// CreateEndpoint creates a new endpoint instance.
func CreateEndpoint(readyChecks ...*abool.AtomicBool) *Endpoint {
	return &Endpoint{
		readyChecks: readyChecks,
	}
}

// AddChecks allows readiness checks to be added to the endpoint.
func (endpoint *Endpoint) AddChecks(readyChecks ...*abool.AtomicBool) {
	endpoint.readyChecks = append(endpoint.readyChecks, readyChecks...)
}

// Name returns the unique ID for this endpoint.
func (endpoint *Endpoint) Name() string {
	return namespace
}

// AddRoutes adds the meta endpoints to the router group.
func (endpoint *Endpoint) AddRoutes(router *gin.RouterGroup) {
	router.GET("/liveness", endpoint.Liveness)
	router.GET("/readiness", endpoint.Readiness)
	router.GET("/health", endpoint.Health)
	router.GET("/routes", endpoint.Routes)
	router.GET("/build", endpoint.Build)
}

// Liveness is used to simply indicate whether a service is hung or in a bad state.
func (endpoint *Endpoint) Liveness(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// Readiness is used to indicate whether a service is ready to take traffic.
func (endpoint *Endpoint) Readiness(ctx *gin.Context) {
	for _, ready := range endpoint.readyChecks {
		if !ready.IsSet() {
			ctx.Status(http.StatusServiceUnavailable)

			return
		}
	}

	ctx.Status(http.StatusNoContent)
}

// Health is a stub of things to come to allow introspection into external dependency states.
func (endpoint *Endpoint) Health(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

// Routes is used to take the entrypoint's base router and print out the routes possible with the server.
func (endpoint *Endpoint) Routes(ctx *gin.Context) {
	router, ok := ctx.MustGet("router").(*gin.Engine)
	if !ok {
		panic(errors.New("failed to convert router from context into `*gin.Engine`"))
	}

	routes := router.Routes()
	routeMap := make(map[string][]string)

	for _, route := range routes {
		if _, ok := routeMap[route.Path]; !ok {
			routeMap[route.Path] = make([]string, 0)
		}

		routeMap[route.Path] = append(routeMap[route.Path], route.Method)
	}

	ctx.JSON(http.StatusOK, routeMap)
}

// Build is used to return any build metadata.
func (endpoint *Endpoint) Build(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, struct {
		Build        string `json:"build,omitempty"`
		Commit       string `json:"commit,omitempty"`
		Date         string `json:"date,omitempty"`
		Version      string `json:"version,omitempty"`
		Program      string `json:"program,omitempty"`
		OS           string `json:"os,omitempty"`
		Architecture string `json:"architecture,omitempty"`
		Arm          string `json:"arm,omitempty"`
	}{
		Build:        build.Build,
		Commit:       build.Commit,
		Date:         build.Date,
		Version:      build.Version,
		Program:      build.Program,
		OS:           build.OS,
		Architecture: build.Architecture,
		Arm:          build.ARM,
	})
}
