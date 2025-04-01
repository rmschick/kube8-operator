package instrumentation

import (
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/runmetrics"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// InstrumentRuntime enables the default runtime metrics.
func InstrumentRuntime() error {
	// Enable default golang runtime metrics
	return errors.Wrap(runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
	}), "failed to enable runtime metrics")
}

// RegisterHTTPViews enables the ochttp client views.
func RegisterHTTPViews() error {
	clientSent := ochttp.ClientSentBytesDistribution
	clientReceived := ochttp.ClientReceivedBytesDistribution
	clientLatency := ochttp.ClientRoundtripLatencyDistribution
	clientCompleted := ochttp.ClientCompletedCount

	clientTags := []tag.Key{ochttp.KeyClientHost, ochttp.KeyClientMethod, ochttp.KeyClientStatus}

	clientSent.TagKeys = clientTags
	clientReceived.TagKeys = clientTags
	clientLatency.TagKeys = clientTags
	clientCompleted.TagKeys = clientTags

	return errors.Wrap(view.Register(
		clientSent,
		clientReceived,
		clientLatency,
		clientCompleted,
	), "failed to register HTTP client metric views")
}

// RegisterHTTPServerViews enables the ochttp server views.
func RegisterHTTPServerViews() error {
	serverRequestCount := ochttp.ServerRequestCountView
	serverResponseBytesView := ochttp.ServerResponseBytesView
	serverRequestCountByMethod := ochttp.ServerRequestCountByMethod
	serverResponseCountByStatusCode := ochttp.ServerResponseCountByStatusCode

	serverTags := []tag.Key{PathTag, HostTag, MethodTag}

	serverRequestCount.TagKeys = serverTags
	serverResponseBytesView.TagKeys = serverTags
	serverRequestCountByMethod.TagKeys = serverTags
	serverResponseCountByStatusCode.TagKeys = serverTags

	return errors.Wrap(view.Register(
		serverRequestCount,
		serverResponseBytesView,
		serverRequestCountByMethod,
		serverResponseCountByStatusCode,
	), "failed to register HTTP server metric views")
}

// InstrumentHTTPClient wrap the http client transport with opencensus's http plugin.
func InstrumentHTTPClient(client *http.Client) *http.Client {
	client.Transport = &ochttp.Transport{Base: client.Transport}

	return client
}

// InstrumentResty sets the resty transport to opencensus's http plugin.
func InstrumentResty(resty *resty.Client) *resty.Client {
	return resty.SetTransport(&ochttp.Transport{})
}

// InstrumentPubSub sets the default pubsub opencensus views.
func InstrumentPubSub() error {
	// https://godoc.org/cloud.google.com/go/pubsub#pkg-variables
	return errors.Wrap(view.Register(
		pubsub.PublishedMessagesView,
		pubsub.PublishLatencyView,
		pubsub.PullCountView,
		pubsub.AckCountView,
		pubsub.NackCountView,
		pubsub.ModAckCountView,
		pubsub.ModAckTimeoutCountView,
		pubsub.StreamOpenCountView,
		pubsub.StreamRetryCountView,
		pubsub.StreamRequestCountView,
		pubsub.StreamResponseCountView,
		ocgrpc.ClientSentBytesPerRPCView,
		ocgrpc.ClientReceivedBytesPerRPCView,
		ocgrpc.ClientRoundtripLatencyView,
		ocgrpc.ClientCompletedRPCsView,
	), "failed to register gRPC and Pubsub metric views")
}
