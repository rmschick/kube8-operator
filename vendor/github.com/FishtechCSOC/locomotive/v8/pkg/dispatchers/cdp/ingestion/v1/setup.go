package ingestion

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	fishtechErrors "github.com/FishtechCSOC/common-go/pkg/errors/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

const (
	customerIDLabelKey = "client-id"
	bucketNameFormat   = "cyderes-uap-%s-%s"
)

func formatBucketName(metadata types.Metadata) string {
	return fmt.Sprintf(bucketNameFormat, metadata.Customer.Name, metadata.Environment)
}

func BucketFromMetadata(client *storage.Client, metadata types.Metadata) *storage.BucketHandle {
	return client.Bucket(formatBucketName(metadata))
}

func BucketWithExistenceCheck(ctx context.Context, client *storage.Client, metadata types.Metadata) (*storage.BucketHandle, error) {
	bucket := BucketFromMetadata(client, metadata)

	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		bucketMeta := map[string]any{
			"bucketName": formatBucketName(metadata),
		}

		switch {
		case errors.Is(err, storage.ErrBucketNotExist):
			return nil, fishtechErrors.CreateAttributeError(err, "bucket does not exist or you do not have access to it", bucketMeta)
		default:
			return nil, fishtechErrors.CreateAttributeError(err, "an error occurred while fetching attrs for the given bucket", bucketMeta)
		}
	}

	if err = validateCustomer(metadata.Customer, attrs); err != nil {
		return nil, fishtechErrors.CreateAttributeError(err, "failed to validate customer", map[string]any{
			"metadataCustomerID":   metadata.Customer.ID,
			"metadataCustomerName": metadata.Customer.Name,
			"bucketName":           formatBucketName(metadata),
			"bucketCustomerID":     attrs.Labels[customerIDLabelKey],
		})
	}

	return bucket, nil
}

func SetupBucketHandle(ctx context.Context, configuration Configuration, metadata types.Metadata, logger *logrus.Entry) *Dispatcher {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucket, err := BucketWithExistenceCheck(ctx, storageClient, metadata)
	if err != nil {
		panic(err)
	}

	return CreateDispatcher(configuration, bucket, logger)
}

// Ensure metadata and bucket labels match the same customer id.
func validateCustomer(customer types.Customer, attrs *storage.BucketAttrs) error {
	if attrs.Labels[customerIDLabelKey] == "" {
		return errors.New("gcp bucket does not have a customer id set")
	}

	if customer.ID != attrs.Labels[customerIDLabelKey] {
		return errors.New("customer id in metadata does not match bucket customer id")
	}

	return nil
}
