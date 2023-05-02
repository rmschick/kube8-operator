package metrics

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// RegisterView sets the default opencensus views.
func RegisterView(tagKeys ...tag.Key) error {
	CyderesTimestampCursorGaugeView.TagKeys = tagKeys

	return view.Register(
		CyderesTimestampCursorGaugeView,
	)
}
