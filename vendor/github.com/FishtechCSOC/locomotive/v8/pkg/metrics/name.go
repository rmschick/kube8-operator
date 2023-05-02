package metrics

import (
	"strings"
)

const (
	Separator = "/"

	MetricPrefix             = "cyderes"
	EntryCountSuffix         = "entry_count"
	EntriesSizeSuffix        = "batch_size"
	RoundtripLatencySuffix   = "roundtrip_latency"
	CompletedCountSuffix     = "completed_count"
	CrossContaminationSuffix = "cross_contamination_detection"
)

func NameBuilder(parts ...string) string {
	return strings.Join(parts, Separator)
}
