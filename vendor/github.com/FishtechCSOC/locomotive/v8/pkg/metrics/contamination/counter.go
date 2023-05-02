package contamination

import (
	"context"
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

type Counter struct {
	metadata           types.Metadata
	contaminationCount int
	endOnce            sync.Once
}

func CreateCounter(metadata types.Metadata) *Counter {
	return &Counter{
		metadata:           metadata,
		contaminationCount: 0,
	}
}

func (counter *Counter) SetContaminationCount(count int) {
	counter.contaminationCount = count
}

func (counter *Counter) End(ctx context.Context) {
	counter.endOnce.Do(func() {
		measurements := []stats.Measurement{
			CyderesCrossContaminationDetectionCount.M(int64(counter.contaminationCount)),
		}

		_ = stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(metrics.CustomerNameTag, counter.metadata.Customer.Name),
			tag.Upsert(metrics.CustomerIDTag, counter.metadata.Customer.ID),
			tag.Upsert(metrics.DataTypeTag, string(counter.metadata.DataType)),
			tag.Upsert(metrics.LogTypeTag, string(counter.metadata.DataType)),
		}, measurements...)
	})
}
