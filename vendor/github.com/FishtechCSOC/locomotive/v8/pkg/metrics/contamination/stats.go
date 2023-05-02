package contamination

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics"
)

// nolint: gochecknoglobals
var (
	CyderesCrossContaminationDetectionCount = stats.Int64(
		metrics.NameBuilder(metrics.MetricPrefix, metrics.CrossContaminationSuffix),
		"Total cross contamination detections",
		stats.UnitDimensionless,
	)
)

// nolint: gochecknoglobals
var (
	CyderesCrossContaminationDetectionView = &view.View{
		Name:        metrics.NameBuilder(metrics.MetricPrefix, metrics.CrossContaminationSuffix),
		Measure:     CyderesCrossContaminationDetectionCount,
		Aggregation: view.Sum(),
		Description: "Count of cross contamination detections, by clientName, clientID, and dataType",
		TagKeys:     []tag.Key{metrics.CustomerNameTag, metrics.CustomerIDTag, metrics.DataTypeTag, metrics.LogTypeTag},
	}
)
