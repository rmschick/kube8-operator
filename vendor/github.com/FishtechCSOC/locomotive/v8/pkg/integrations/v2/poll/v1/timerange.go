package poll

import (
	"time"
)

// TimeRange is a data structure that contains a start and end time which does not enforce start to always be less than
// the end time.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Difference returns the difference in duration between start and end time.
func (timeRange TimeRange) Difference() time.Duration {
	return timeRange.End.Sub(timeRange.Start)
}
