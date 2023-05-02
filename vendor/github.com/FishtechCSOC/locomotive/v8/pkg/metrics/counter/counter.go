package counter

import (
	"context"
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics"
	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics/tracker"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

type Counter struct {
	entries *types.LogEntries
	status  types.Status
	endOnce sync.Once
}

// CreateCounter creates a instance of counter.
func CreateCounter(metadata types.Metadata) *Counter {
	return &Counter{
		entries: &types.LogEntries{
			Metadata: metadata,
		},
		status: types.Failure,
	}
}

func (counter *Counter) SetEntries(entries *types.LogEntries) {
	if entries == nil {
		return
	}

	counter.entries = entries
}

// SetStatus sets the tracker to the given status.
func (counter *Counter) SetStatus(status types.Status) {
	counter.status = status
}

// End stops the tracker and records the metrics.
func (counter *Counter) End(ctx context.Context) {
	counter.endOnce.Do(func() {
		entryCount := counter.entries.EntryCount()
		batchSize := counter.entries.SizeInBytes()

		measurements := []stats.Measurement{
			tracker.CyderesEntriesCount.M(int64(entryCount)),
			tracker.CyderesBatchSize.M(int64(batchSize)),
		}

		_ = stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(metrics.CustomerNameTag, counter.entries.Metadata.Customer.Name),
			tag.Upsert(metrics.CustomerIDTag, counter.entries.Metadata.Customer.ID),
			tag.Upsert(metrics.DataTypeTag, string(counter.entries.Metadata.DataType)),
			tag.Upsert(metrics.LogTypeTag, string(counter.entries.Metadata.DataType)),
			tag.Upsert(metrics.StatusTag, string(counter.status)),
		}, measurements...)
	})
}
