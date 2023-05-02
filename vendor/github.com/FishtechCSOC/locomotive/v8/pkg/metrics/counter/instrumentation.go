package counter

import (
	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics/tracker"
)

// RegisterView sets the default opencensus views.
func RegisterView() error {
	return tracker.RegisterView()
}
