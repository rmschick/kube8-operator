package ingestion

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/math"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

const (
	dispatcherType = "cdp-ingestion"

	timeFormat = "2006/01/02/15/04"
)

var _ integrations.Dispatcher = (*Dispatcher)(nil)

// Dispatcher is the object that handles the business logic of dispatching entries to Azure Log Analytics.
type Dispatcher struct {
	gcsBucket     *storage.BucketHandle
	configuration Configuration
	logger        *logrus.Entry
}

// CreateDispatcher creates a dispatcher instance.
func CreateDispatcher(configuration Configuration, gcsBucket *storage.BucketHandle, logger *logrus.Entry) *Dispatcher {
	dispatcher := &Dispatcher{
		gcsBucket:     gcsBucket,
		configuration: configuration,
	}

	dispatcher.logger = integrations.SetupDispatcherLogger(dispatcherType, dispatcher, logger)

	return dispatcher
}

// Type returns a unique ID of the dispatcher type.
func (dispatcher *Dispatcher) Type() string {
	return dispatcherType
}

// Dispatch accepts an LogEntries object to dispatch and returns whether it was successful or not.
func (dispatcher *Dispatcher) Dispatch(ctx context.Context, entries *types.LogEntries) types.Status {
	if len(entries.Entries) <= 0 {
		dispatcher.logger.Info("No entries to dispatch, sending success")

		return types.Success
	}

	objectPath := fmt.Sprintf("%s/%s/%s/%s-%s.gz", build.Program, string(entries.Metadata.DataType), time.Now().UTC().Format(timeFormat), strconv.FormatInt(time.Now().UnixNano(), math.DecimalNumeralSystem), entries.Metadata.Instance)
	obj := dispatcher.gcsBucket.Object(objectPath)
	writer := obj.NewWriter(ctx)
	writer.Metadata = entries.Metadata.MarshalToMetadata()
	writer.ContentEncoding = "application/gzip"
	writer.ContentType = "text/plain"
	writer.Metadata["content-encoding"] = "application/gzip"
	writer.Metadata["content-type"] = "text/plain"

	var buf bytes.Buffer

	err := dispatcher.compressLogEntries(&buf, entries)
	if err != nil {
		dispatcher.logger.WithError(err).Error("Failed to compress log entries")
	}

	_, err = writer.Write(buf.Bytes())
	if err != nil {
		dispatcher.logger.WithError(err).Error("Failed to write data to storage writer")

		return types.Failure
	}

	err = writer.Close()
	if err != nil {
		dispatcher.logger.WithError(err).Error("Failed to write data to GCS.")

		return types.Failure
	}

	dispatcher.logger.Debug("Successfully posted message to GCS.")

	return types.Success
}

func (dispatcher *Dispatcher) compressLogEntries(buf io.Writer, entries *types.LogEntries) error {
	zipWriter := gzip.NewWriter(buf)

	for _, entry := range entries.Entries {
		var rawLog []byte

		switch {
		case entry.Log != "":
			rawLog = []byte(entry.Log + "\n")
		default:
			dispatcher.logger.Debug("No useable data in entry, continuing")

			continue
		}

		_, err := zipWriter.Write(rawLog)
		if err != nil {
			dispatcher.logger.WithError(err).Error("Failed to write to buffer, continuing")

			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("error while closing GZip writer: %w", err)
	}

	return nil
}
