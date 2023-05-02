package runner

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/metrics/counter"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

type Worker struct {
	name          string
	configuration Configuration
	dispatcher    integrations.Dispatcher
	buffer        []types.LogEntries
	stream        <-chan *types.LogEntries
	logger        *logrus.Entry
	totalEntries  int
	totalSize     int
}

func CreateWorker(name string, configuration Configuration, dispatcher integrations.Dispatcher, stream <-chan *types.LogEntries, logger *logrus.Entry) *Worker {
	worker := &Worker{
		name:          name,
		configuration: configuration,
		dispatcher:    dispatcher,
		stream:        stream,
		logger:        logger.WithField("worker", name),
		buffer:        make([]types.LogEntries, 0),
	}

	return worker
}

// Run is the main entrypoint for beginning processing work. Run blocks until either the context.Context given is
// closed or the stream this was created with is closed and empty.
func (worker *Worker) Run(ctx context.Context, done func(), _ bool) {
	// nolint: contextcheck
	defer worker.close(done)

	var (
		tick   <-chan time.Time
		ticker *time.Ticker
	)

	switch {
	case worker.configuration.Batch.WaitTime <= 0:
		tick = make(chan time.Time)
	default:
		ticker = time.NewTicker(worker.configuration.Batch.WaitTime)
		tick = ticker.C
	}

	for {
		select {
		case <-ctx.Done():
			worker.logger.WithField("reason", ctx.Err()).Info("Context closed, exiting worker")

			return
		case timestamp, ok := <-tick:
			if !ok {
				continue
			}

			worker.logger.WithField("time", timestamp).Debug("Timer expired, flushing")
			worker.flush(ctx, ticker)
		case entries, ok := <-worker.stream:
			if !ok {
				return
			}

			if entries == nil {
				continue
			}

			worker.append(ctx, ticker, entries)
		}
	}
}

// nolint: unused
// shouldClose is special logic to ensure that drain does not actually exit without processing what is left on the
// stream during the parent Runner's drain. May actually be superfluous.
func (worker *Worker) shouldClose(ok, drain, empty bool) bool {
	if !(ok || drain) {
		return true
	}

	if !ok && drain && empty {
		return true
	}

	return false
}

// flush handles both pushing data and ensuring that the buffer and all accompanying variables are appropriately reset.
// If the buffer is empty when flush is called, a guard is in place to ensure flush is a no-op (no reset occurs).
func (worker *Worker) flush(ctx context.Context, ticker *time.Ticker) {
	if len(worker.buffer) <= 0 {
		return
	}

	statusFuncs := make([]func(types.Status), 0, len(worker.buffer))
	allEntries := &types.LogEntries{
		Metadata: worker.buffer[0].Metadata,
		Entries:  make([]types.LogEntry, 0, worker.totalEntries),
	}

	for _, entries := range worker.buffer {
		allEntries.Append(entries.Entries...)

		if entries.Status != nil {
			statusFuncs = append(statusFuncs, entries.Status)
		}
	}

	allEntries.Status = func(status types.Status) {
		for _, f := range statusFuncs {
			f(status)
		}
	}

	// We may pull a fluentd and configure behavior here, whether we just drop or block and retry until this works
	status := worker.processEntries(ctx, allEntries)

	if status == types.Failure {
		worker.logger.Info("Failed to dispatch batch, dropping entries")
	}

	worker.reset(ticker)
}

// append appends given entries and then flushes if buffer is considered full. If all Batch configuration is
// non-positive, the buffer is always flushed with any new entry to simulate as if there was no batching.
func (worker *Worker) append(ctx context.Context, ticker *time.Ticker, entries *types.LogEntries) {
	worker.totalEntries += entries.EntryCount()
	worker.totalSize += entries.SizeInBytes()
	worker.buffer = append(worker.buffer, *entries)

	switch {
	case worker.configuration.Batch.SizeInBytes > 0 && worker.totalSize > worker.configuration.Batch.SizeInBytes:
		worker.logger.WithField("sizeInBytes", worker.totalSize).Debug("Byte size limit reached, flushing")
		worker.flush(ctx, ticker)
	case worker.configuration.Batch.EntryCount > 0 && worker.totalEntries > worker.configuration.Batch.EntryCount:
		worker.logger.WithField("entryCount", worker.totalEntries).Debug("Entry count limit reached, flushing")
		worker.flush(ctx, ticker)
	case worker.configuration.Batch.WaitTime <= 0 &&
		worker.configuration.Batch.SizeInBytes <= 0 &&
		worker.configuration.Batch.EntryCount <= 0:
		worker.flush(ctx, ticker)
	}
}

// reset creates a new buffer using the capacity of the former buffer and resets the totals used for determining when
// to flush the buffer. The ticker is also reset in case it was not responsible for the flush/reset and full duration
// is waited.
func (worker *Worker) reset(ticker *time.Ticker) {
	worker.buffer = make([]types.LogEntries, 0, cap(worker.buffer))
	worker.totalSize = 0
	worker.totalEntries = 0

	if ticker != nil && worker.configuration.Batch.WaitTime > 0 {
		ticker.Reset(worker.configuration.Batch.WaitTime)
	}
}

// processEntries is responsible for both recording metrics around entries as well as dispatching them. Empty dataTypes
// are treated as failures while empty entries are skipped and treated as successes. True or false is returned based on
// whether or not the status of the entries was successful or not.
func (worker *Worker) processEntries(ctx context.Context, entries *types.LogEntries) types.Status {
	entriesCounter := counter.CreateCounter(types.CreateMetadata())

	defer entriesCounter.End(ctx)

	entriesCounter.SetEntries(entries)

	var status types.Status

	switch {
	case len(entries.Entries) <= 0:
		worker.logger.Info("No entries to push, continuing stream")

		status = types.Success
	case entries.Metadata.DataType == "":
		worker.logger.Info("Empty data type detected, cannot push data without a data type")

		status = types.Failure
	default:
		status = worker.dispatcher.Dispatch(ctx, entries)
	}

	entriesCounter.SetStatus(status)
	entries.ReportStatus(status)

	return status
}

// close ensures that everything is cleaned up properly including flushing whatever happened to be on the buffer after
// a cleanup was triggered but before a flush condition was met.
func (worker *Worker) close(done func()) {
	worker.logger.Info("Attempting to flush buffer before closing")

	ctx, cancel := context.WithTimeout(context.Background(), worker.configuration.DrainTimeout)

	worker.flush(ctx, nil)
	cancel()
	done()
}
