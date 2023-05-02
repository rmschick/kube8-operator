package setup

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/FishtechCSOC/common-go/pkg/web/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/endpoints/debug/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/endpoints/meta/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/endpoints/metrics/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/server/v1"
	"github.com/sirupsen/logrus"
)

// BuildHTTPServer creates a server with prometheus metrics from a config, with a readiness check.
func BuildHTTPServer(configuration server.Configuration, prometheusExporter *prometheus.Exporter, logger *logrus.Entry) *server.Server {
	metricsEndpoint := metrics.CreateEndpoint(prometheusExporter)
	metaEndpoint := meta.CreateEndpoint()
	debugEndpoint := debug.CreateEndpoint()

	httpServer, err := server.CreateServer([]web.Endpoint{metricsEndpoint, metaEndpoint, debugEndpoint}, []web.Middleware{}, configuration, logger)
	if err != nil {
		panic(err)
	}

	metaEndpoint.AddChecks(httpServer.Ready())

	return httpServer
}
