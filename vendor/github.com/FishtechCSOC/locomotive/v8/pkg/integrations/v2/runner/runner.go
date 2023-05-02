package runner

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/runner/chunk"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

// Runner handles the management of the actual lifecycle of integrations as well as common concerns.
type Runner struct {
	name          string
	retriever     integrations.Retriever
	dispatcher    integrations.Dispatcher
	chunker       chunk.Chunker
	stream        chan *types.LogEntries
	workers       chan *types.LogEntries
	configuration Configuration
	done          chan bool
	logger        *logrus.Entry
}

// CreateRunner sets up a new Runner instance.
func CreateRunner(name string, configuration Configuration, retriever integrations.Retriever, dispatcher integrations.Dispatcher, chunker chunk.Chunker, logger *logrus.Entry) *Runner {
	runner := &Runner{
		name:          name,
		retriever:     retriever,
		dispatcher:    dispatcher,
		chunker:       chunker,
		stream:        make(chan *types.LogEntries, configuration.WorkerCount),
		workers:       make(chan *types.LogEntries, configuration.WorkerCount),
		configuration: configuration,
		done:          make(chan bool, 1),
	}

	runner.logger = logging.CreateTypeLogger(logger, name, runner)

	return runner
}

// Run handles creates workers to send types.LogEntries to be dispatched but also is responsible for starting the
// integrations.Retriever and the chunking logic within the runner. Run blocks until the context is closed and all
// internal workers exit safely.
func (runner *Runner) Run(ctx context.Context, cancel context.CancelFunc) {
	workers := make([]*Worker, runner.configuration.WorkerCount)
	workerGroup := &sync.WaitGroup{}

	for i := range workers {
		workers[i] = CreateWorker(strconv.Itoa(i), runner.configuration, runner.dispatcher, runner.workers, runner.logger)
	}

	// nolint: contextcheck
	defer runner.close(cancel)

	if runner.configuration.WorkerCount < 0 {
		runner.logger.WithField("workerCount", runner.configuration.WorkerCount).Info("Non-positive worker totalEntries, exiting")

		return
	}

	runner.logger.Info("Starting streaming")

	go runner.retriever.Retrieve(ctx, cancel, runner.stream)
	go runner.chunker.Chunk(ctx, runner.stream, runner.workers)

	for _, worker := range workers {
		workerGroup.Add(1)

		go worker.Run(ctx, workerGroup.Done, false)
	}

	workerGroup.Wait()
	runner.logger.WithField("reason", ctx.Err()).Info("runner closed due to context closed")
}

// Wait blocks until the Runner is successfully cleaned up.
func (runner *Runner) Wait() {
	runner.logger.Info("Began waiting for retriever to stop")
	<-runner.done
	runner.logger.Info("Done waiting for retriever to stop")
}

// drain in case one or both of the channels for types.LogEntries are not empty, we double check here and try to empty
// them.
func (runner *Runner) drain() {
	if runner.configuration.DrainTimeout <= 0 {
		runner.logger.Info("Retriever was not configured to drain, closing immediately")

		return
	}

	if len(runner.stream) <= 0 && len(runner.workers) <= 0 {
		runner.logger.Info("Retriever stream was empty so no need to drain")

		return
	}

	runner.logger.Info("Draining retriever of all data before closing")

	workers := make([]*Worker, runner.configuration.WorkerCount)
	workerGroup := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), runner.configuration.DrainTimeout)

	go runner.chunker.Chunk(ctx, runner.stream, runner.workers)

	for i := range workers {
		workers[i] = CreateWorker(strconv.Itoa(i), runner.configuration, runner.dispatcher, runner.workers, runner.logger)

		workerGroup.Add(1)

		go workers[i].Run(ctx, workerGroup.Done, true)
	}

	go func(workerGroup *sync.WaitGroup, cancel context.CancelFunc) {
		workerGroup.Wait()
		cancel()
	}(workerGroup, cancel)

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		runner.logger.Info("Timed out waiting for retriever to drain, closing")
	}
}

// close is responsible for triggering all the cleanup necessary to exit the runner.
func (runner *Runner) close(cancel context.CancelFunc) {
	runner.logger.Info("Closing runner")
	cancel()
	runner.retriever.Wait()
	close(runner.stream)
	// nolint: contextcheck
	runner.drain()
	close(runner.done)
}
