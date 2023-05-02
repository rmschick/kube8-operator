package setup

import (
	"github.com/FishtechCSOC/common-go/pkg/metrics/instrumentation"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics/contamination"
	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics/counter"
)

// RegisterMetricViews attempts to register runtime, http, and poller metrics.
func RegisterMetricViews(logger *logrus.Entry) {
	err := instrumentation.InstrumentRuntime()
	if err != nil {
		logger.WithError(err).Info("Failed to instrument go runtime metrics")
	}

	err = instrumentation.RegisterHTTPViews()
	if err != nil {
		logger.WithError(err).Info("Failed to register http client metrics")
	}

	err = counter.RegisterView()
	if err != nil {
		logger.WithError(err).Info("Failed to register poller metrics")
	}

	err = contamination.RegisterView()
	if err != nil {
		logger.WithError(err).Info("Failed to register contamination metrics")
	}
}
