package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics"
)

const CursorSuffix = "cursor_timestamp"

// nolint: gochecknoglobals
var (
	CyderesTimestampCursorGauge = stats.Int64(
		metrics.NameBuilder(metrics.MetricPrefix, CursorSuffix),
		"Unix Timestamp of time related cursors",
		stats.UnitSeconds,
	)
)

// nolint: gochecknoglobals
var (
	CyderesTimestampCursorGaugeView = &view.View{
		Name:        metrics.NameBuilder(metrics.MetricPrefix, CursorSuffix),
		Measure:     CyderesTimestampCursorGauge,
		Aggregation: view.LastValue(),
		Description: "Timestamp of the cursor's position at last poll",
	}
)
