package metrics

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/gin-gonic/gin"
)

const (
	namespace = "metrics"
)

// Endpoint is the manager around the endpoint routes.
type Endpoint struct {
	prometheusExporter *prometheus.Exporter
}

// CreateEndpoint creates a new endpoint instance.
func CreateEndpoint(prometheusExporter *prometheus.Exporter) *Endpoint {
	return &Endpoint{
		prometheusExporter: prometheusExporter,
	}
}

// Name returns the unique ID for this endpoint.
func (endpoint *Endpoint) Name() string {
	return namespace
}

// AddRoutes adds the meta endpoints to the router group.
func (endpoint *Endpoint) AddRoutes(router *gin.RouterGroup) {
	if endpoint.prometheusExporter != nil {
		router.GET("prometheus", endpoint.PrometheusMetrics)
	}
}

// PrometheusMetrics outputs prometheus formatted metrics to the response.
func (endpoint *Endpoint) PrometheusMetrics(ctx *gin.Context) {
	endpoint.prometheusExporter.ServeHTTP(ctx.Writer, ctx.Request)
}
