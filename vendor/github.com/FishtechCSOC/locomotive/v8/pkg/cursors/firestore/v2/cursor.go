package firestore

import (
	"context"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
)

var _ integrations.Cursor = (*Cursor)(nil)

type Cursor struct {
	config   Configuration
	document CursorDocument
	client   *firestore.DocumentRef
	logger   *logrus.Entry
	mutex    sync.Mutex
}

func CreateCursor(ctx context.Context, config Configuration, client *firestore.DocumentRef, logger *logrus.Entry) (*Cursor, error) {
	cursor := &Cursor{
		config: config,
		client: client,
	}

	cursor.logger = integrations.SetupCursorLogger(CursorType, cursor, logger)

	if _, err := cursor.Load(ctx); err != nil {
		return nil, err
	}

	return cursor, nil
}

func (cursor *Cursor) Store(ctx context.Context, value string) error {
	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	cursor.document.Value = value

	return cursor.syncDocument(ctx)
}

func (cursor *Cursor) Load(ctx context.Context) (string, error) {
	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	snapshot, err := cursor.client.Get(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to load document from firestore")
	}

	err = snapshot.DataTo(&cursor.document)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse document from firestore")
	}

	return cursor.document.Value, nil
}

func (cursor *Cursor) Test(ctx context.Context) error {
	if _, err := cursor.client.Get(ctx); err != nil {
		return errors.Wrap(err, "failed to load document from firestore")
	}

	return nil
}

func (cursor *Cursor) Close(ctx context.Context) {
	cursor.mutex.Lock()
	defer cursor.mutex.Unlock()

	if err := cursor.syncDocument(ctx); err != nil {
		cursor.logger.WithError(err).Error("failed to sync document on cursor close")
	}
}

func (cursor *Cursor) syncDocument(ctx context.Context) error {
	cursor.document.LastUpdate = time.Now().UTC()
	cursor.document.Expiration = cursor.document.LastUpdate.Add(cursor.config.TimeToLive)

	_, err := cursor.client.Update(ctx, []firestore.Update{
		{
			Path:  "value",
			Value: cursor.document.Value,
		},
		{
			Path:  "lastUpdate",
			Value: cursor.document.LastUpdate,
		},
		{
			Path:  "expiration",
			Value: cursor.document.Expiration,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to sync document to firestore")
	}

	return nil
}
