package tracker

import (
	"context"
	"sync"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

type Tracker struct {
	entries *types.LogEntries
	status  types.Status
	endOnce sync.Once
}

// CreateTracker creates a instance of tracker.
func CreateTracker(metadata types.Metadata) *Tracker {
	return &Tracker{
		entries: &types.LogEntries{
			Metadata: metadata,
		},
		status: types.Failure,
	}
}

// SetSuccess sets the tracker as being successful.
func (tracker *Tracker) SetEntries(entries *types.LogEntries) {
	if entries == nil {
		return
	}

	tracker.entries = entries
}

// SetStatus sets the tracker to the given status.
func (tracker *Tracker) SetStatus(status types.Status) {
	tracker.status = status
}

// End stops the tracker and records the metrics.
func (tracker *Tracker) End(ctx context.Context) {
	tracker.endOnce.Do(func() {
	})
}
