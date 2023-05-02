package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

func SetupCursor(ctx context.Context, metadata types.Metadata, configuration Configuration, shard string, logger *logrus.Entry) *Cursor {
	client, err := firestore.NewClient(ctx, configuration.ProjectID)
	if err != nil {
		panic(err)
	}

	currentTime := time.Now().UTC()

	documentRef, err := GetOrCreateDocument(ctx, client, metadata, shard, currentTime, currentTime.Add(configuration.TimeToLive))
	if err != nil {
		panic(err)
	}

	cursor, err := CreateCursor(ctx, configuration, documentRef, logger)
	if err != nil {
		panic(err)
	}

	return cursor
}

func GetOrCreateDocument(ctx context.Context, client *firestore.Client, metadata types.Metadata, shard string, currentTime, expiration time.Time) (*firestore.DocumentRef, error) {
	cursorDocument := CursorDocument{
		Collector:  build.Program,
		Customer:   metadata.Customer.ID,
		Instance:   metadata.Instance,
		Shard:      shard,
		Value:      "",
		LastUpdate: currentTime,
		Expiration: expiration,
	}
	collectionRef := client.Collection(v1CollectionRef)
	documentIterator := collectionRef.
		Where("collector", "==", cursorDocument.Collector).
		Where("customer", "==", cursorDocument.Customer).
		Where("instance", "==", cursorDocument.Instance).
		Where("shard", "==", cursorDocument.Shard).
		Documents(ctx)

	documents, err := documentIterator.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while iterating collection: %w", err)
	}

	if len(documents) <= 0 {
		return createDocument(ctx, collectionRef, cursorDocument)
	}

	if len(documents) > 1 {
		return nil, MultipleCursorFoundError
	}

	return documents[0].Ref, nil
}

func createDocument(ctx context.Context, collectionRef *firestore.CollectionRef, cursorDocument CursorDocument) (*firestore.DocumentRef, error) {
	doc, _, err := collectionRef.Add(ctx, cursorDocument)
	if err != nil {
		return nil, CreateCursorDocumentError
	}

	if doc == nil {
		return nil, NoCursorFoundError
	}

	return doc, nil
}
