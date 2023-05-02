package linear

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/runner/chunk"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

var _ chunk.Chunker = (*Chunker)(nil)

type Chunker struct {
	configuration chunk.Configuration
	logger        *logrus.Entry
}

func CreateChunker(configuration chunk.Configuration, logger *logrus.Entry) *Chunker {
	chunker := &Chunker{
		configuration: configuration,
	}

	chunker.logger = chunk.SetupChunkerLogger("linear", chunker, logger)

	return chunker
}

func (chunker *Chunker) Chunk(ctx context.Context, input <-chan *types.LogEntries, output chan<- *types.LogEntries) {
	for {
		select {
		case <-ctx.Done():
			chunker.logger.WithField("reason", ctx.Err()).Info("Context closed, exiting chunker")

			return
		case entries, ok := <-input:
			if !ok {
				chunker.logger.Info("Input channel closed, exiting chunker")

				return
			}

			chunker.chunk(ctx, entries, output)
		}
	}
}

func (chunker *Chunker) chunk(ctx context.Context, entries *types.LogEntries, output chan<- *types.LogEntries) {
	switch {
	case chunker.configuration.SizeInBytes > 0:
		chunker.chunkBySize(ctx, entries, output)
	case chunker.configuration.EntryCount > 0:
		chunker.chunkByCount(ctx, entries, output)
	default:
		if err := integrations.SendToStream(ctx, output, entries); err != nil {
			chunker.logger.WithError(err).Error("Failed to push chunk to stream")

			return
		}
	}
}

func (chunker *Chunker) chunkBySize(ctx context.Context, entries *types.LogEntries, output chan<- *types.LogEntries) {
	cursor := 0
	chunkSize := 0

	for index, entry := range entries.Entries {
		entrySize := len(entry.Log)

		if chunkSize+entrySize > chunker.configuration.SizeInBytes {
			if err := integrations.SendToStream(ctx, output, &types.LogEntries{
				Metadata: entries.Metadata,
				Entries:  entries.Entries[cursor:index],
				Status:   entries.Status,
			}); err != nil {
				chunker.logger.WithField("reason", ctx.Err()).Info("Context closed, exiting chunker")

				return
			}

			cursor = index
			chunkSize = 0
		}

		chunkSize += entrySize
	}

	if err := integrations.SendToStream(ctx, output, &types.LogEntries{
		Metadata: entries.Metadata,
		Entries:  entries.Entries[cursor:],
		Status:   entries.Status,
	}); err != nil {
		chunker.logger.WithField("reason", ctx.Err()).Info("Context closed, exiting chunker")

		return
	}
}

func (chunker *Chunker) chunkByCount(ctx context.Context, entries *types.LogEntries, output chan<- *types.LogEntries) {
	cursor := 0
	chunkCount := 0

	for index := range entries.Entries {
		if chunkCount+1 > chunker.configuration.EntryCount {
			if err := integrations.SendToStream(ctx, output, &types.LogEntries{
				Metadata: entries.Metadata,
				Entries:  entries.Entries[cursor:index],
				Status:   entries.Status,
			}); err != nil {
				return
			}

			cursor = index
			chunkCount = 0
		}

		chunkCount++
	}

	if err := integrations.SendToStream(ctx, output, &types.LogEntries{
		Metadata: entries.Metadata,
		Entries:  entries.Entries[cursor:],
		Status:   entries.Status,
	}); err != nil {
		chunker.logger.WithField("reason", ctx.Err()).Info("Context closed, exiting chunker")

		return
	}
}
