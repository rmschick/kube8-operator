package internal

import (
	"context"
	"time"

	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/stream"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
	"github.com/sirupsen/logrus"
)

var _ integrations.Retriever = (*Streamer)(nil)

type Streamer struct {
	configuration stream.Configuration
	cursor        integrations.Cursor
	ticker        *time.Ticker
	logger        *logrus.Entry
	metadata      types.Metadata
	done          chan bool
}

func CreateStreamer(metadata types.Metadata, configuration stream.Configuration, cursor integrations.Cursor, logger *logrus.Entry) *Streamer {
	streamer := &Streamer{
		metadata:      metadata,
		configuration: configuration,
		cursor:        cursor,
		ticker:        time.NewTicker(time.Second),
		done:          make(chan bool, 1),
	}

	streamer.logger = stream.SetupStreamerLogger(string(metadata.DataType), streamer, logger)

	return streamer
}

func (streamer *Streamer) Type() string {
	return build.Program
}

// Retrieve handles the main lifecycle loop using a done context as the closing case.
func (streamer *Streamer) Retrieve(ctx context.Context, cancel context.CancelFunc, stream chan<- *types.LogEntries) {
	defer streamer.close(cancel)

	for {
		select {
		case <-ctx.Done():
			return
		case <-streamer.ticker.C:
			currentTime := time.Now()
			entries := &types.LogEntries{
				// Status: func(status types.Status) {
				// 	streamer.logger.Info(status)
				// },
				Metadata: streamer.metadata,
				Entries: []types.LogEntry{
					{
						Log:       "foo",
						Timestamp: currentTime,
					},
					{
						Log:       "bar",
						Timestamp: currentTime,
					},
				},
			}

			if err := integrations.SendToStream(ctx, stream, entries); err != nil {
				streamer.logger.WithError(err).Info("Failed to push data to stream")

				return
			}
		}
	}
}

// Wait is a helper function that does not return until everything was shutdown.
func (streamer *Streamer) Wait() {
	<-streamer.done
}

func (streamer *Streamer) close(cancel context.CancelFunc) {
	if streamer.ticker != nil {
		streamer.ticker.Stop()
	}

	close(streamer.done)

	cancel()
}
