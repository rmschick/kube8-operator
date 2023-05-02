package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
)

const (
	CursorType = "metrics"
)

var _ integrations.Cursor = (*Cursor)(nil)

type Cursor struct {
	cursor integrations.Cursor
	logger *logrus.Entry
	mutex  sync.Mutex
	tags   []tag.Mutator
}

func CreateCursor(cursor integrations.Cursor, logger *logrus.Entry, tags map[string]string) *Cursor {
	mutators := make([]tag.Mutator, 0, len(tags))
	for k, v := range tags {
		mutators = append(mutators, tag.Upsert(tag.MustNewKey(k), v))
	}

	wrappedCursor := &Cursor{
		cursor: cursor,
		logger: logger,
		tags:   mutators,
	}

	wrappedCursor.logger = integrations.SetupCursorLogger(CursorType, cursor, logger)

	return wrappedCursor
}

// Store stores the value in the wrapped cursor, then records the metric.
func (cursor *Cursor) Store(ctx context.Context, value string) error {
	timestamp, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return errors.Wrap(err, "failed to parse cursor timestamp")
	}

	err = cursor.cursor.Store(ctx, value)
	if err != nil {
		return errors.Wrap(err, "failed to store cursor timestamp")
	}

	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	err = stats.RecordWithTags(ctx,
		cursor.tags,
		CyderesTimestampCursorGauge.M(timestamp.Unix()))

	if err != nil {
		cursor.logger.WithError(err).Error("failed to record cursor metric")
	}

	return nil
}

func (cursor *Cursor) Load(ctx context.Context) (string, error) {
	return cursor.cursor.Load(ctx)
}

func (cursor *Cursor) Test(ctx context.Context) error {
	return cursor.cursor.Test(ctx)
}

func (cursor *Cursor) Close(ctx context.Context) {
	cursor.cursor.Close(ctx)
}
