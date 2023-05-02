package tracker

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics"
)

// The following measures are supported for use in custom views.
// nolint: gochecknoglobals
var (
	CyderesEntriesCount = stats.Int64(
		metrics.NameBuilder(metrics.MetricPrefix, metrics.EntryCountSuffix),
		"Total entries",
		stats.UnitDimensionless,
	)
	CyderesBatchSize = stats.Int64(
		metrics.NameBuilder(metrics.MetricPrefix, metrics.EntriesSizeSuffix),
		"Size of log batch",
		stats.UnitBytes,
	)
)

// You still need to register these views for data to actually be collected.
// nolint: gochecknoglobals
var (
	CyderesEntriesCountView = &view.View{
		Name:        metrics.NameBuilder(metrics.MetricPrefix, metrics.EntryCountSuffix),
		Measure:     CyderesEntriesCount,
		Aggregation: view.Sum(),
		Description: "Count of entries, by clientName, clientID, dataType, and status",
		TagKeys:     []tag.Key{metrics.CustomerNameTag, metrics.CustomerIDTag, metrics.DataTypeTag, metrics.LogTypeTag, metrics.StatusTag},
	}

	CyderesCompletedCountView = &view.View{
		Name:        metrics.NameBuilder(metrics.MetricPrefix, metrics.CompletedCountSuffix),
		Measure:     CyderesEntriesCount,
		Aggregation: view.Count(),
		Description: "Count of completed requests, by clientName, clientID, dataType, and status",
		TagKeys:     []tag.Key{metrics.CustomerNameTag, metrics.CustomerIDTag, metrics.DataTypeTag, metrics.LogTypeTag, metrics.StatusTag},
	}

	CyderesBatchSizeView = &view.View{
		Name:        metrics.NameBuilder(metrics.MetricPrefix, metrics.EntriesSizeSuffix),
		Measure:     CyderesBatchSize,
		Aggregation: view.Sum(),
		Description: "Size of log entries, by clientName, clientID, dataType, and status",
		TagKeys:     []tag.Key{metrics.CustomerNameTag, metrics.CustomerIDTag, metrics.DataTypeTag, metrics.LogTypeTag, metrics.StatusTag},
	}
)
