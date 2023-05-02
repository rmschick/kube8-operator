package contamination

import "go.opencensus.io/stats/view"

// RegisterView sets the default opencensus views.
func RegisterView() error {
	return view.Register(
		CyderesCrossContaminationDetectionView,
	)
}
