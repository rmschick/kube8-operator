package throttled

import (
	"context"
	"sync"
	"time"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
)

var _ integrations.Cursor = (*Cursor)(nil)

type Cursor struct {
	value         string
	configuration Configuration
	cursor        integrations.Cursor
	lastSync      time.Time
	mutex         sync.Mutex
}

func NewCursor(configuration Configuration, cursor integrations.Cursor) (*Cursor, error) {
	newCursor := &Cursor{
		configuration: configuration,
		cursor:        cursor,
	}

	return newCursor, nil
}

func (cursor *Cursor) Store(ctx context.Context, value string) error {
	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	if !isTimeToSync(cursor.lastSync, cursor.configuration.ThrottleOffset) {
		cursor.value = value

		return nil
	}

	if err := cursor.cursor.Store(ctx, value); err != nil {
		return err
	}

	cursor.value = value
	cursor.lastSync = time.Now()

	return nil
}

func (cursor *Cursor) Load(ctx context.Context) (string, error) {
	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	if cursor.value != "" {
		return cursor.value, nil
	}

	value, err := cursor.cursor.Load(ctx)
	if err != nil {
		return "", err
	}

	cursor.value = value
	cursor.lastSync = time.Now()

	return cursor.value, nil
}

func (cursor *Cursor) Test(ctx context.Context) error {
	return cursor.cursor.Test(ctx)
}

func (cursor *Cursor) Close(ctx context.Context) {
	cursor.cursor.Close(ctx)
}

// true if lastSync older than offset minutes (default 5).
func isTimeToSync(lastSync time.Time, offset time.Duration) bool {
	return time.Since(lastSync) > offset
}
