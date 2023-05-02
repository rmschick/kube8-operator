package runnable

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Group is just a helper for managing several Runnable workers and shutting them down when the OS sends specific signals.
type Group struct {
	closed  bool
	mutex   *sync.Mutex
	logger  *logrus.Entry
	stop    chan bool
	signals chan os.Signal
	Runners []Runnable
}

// CreateGroup returns a Group instance with the given runners and a new signal channel.
func CreateGroup(logger *logrus.Entry, runners ...Runnable) *Group {
	runnerGroup := &Group{
		closed:  false,
		mutex:   &sync.Mutex{},
		logger:  logger,
		stop:    make(chan bool, 1),
		signals: make(chan os.Signal, 1),
		Runners: runners,
	}

	signal.Notify(runnerGroup.signals, syscall.SIGINT, syscall.SIGTERM)

	return runnerGroup
}

// Run starts up all Runnable workers and creates a worker that waits to receive a signals from the OS.
func (group *Group) Run(ctx context.Context) {
	signalCtx, signalCancel := context.WithCancel(ctx)
	workerCtx, workerCancel := context.WithCancel(ctx)

	go group.waitForWorkers(workerCtx, workerCancel, signalCancel)
	go group.waitForSignal(workerCancel, signalCancel)

	for _, runner := range group.Runners {
		go runner.Run(signalCtx, workerCancel)
	}
}

// Wait sets up a bunch of workers to wait for each Runnable to be successfully cleaned up before returning.
func (group *Group) Wait() {
	var waitGroup sync.WaitGroup

	for _, runner := range group.Runners {
		waitGroup.Add(1)

		go group.waitForWorker(runner.Wait, &waitGroup)
	}

	waitGroup.Wait()
	group.logger.Info("Done waiting for all workers")
}

func (group *Group) waitForWorker(f func(), waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	f()
	group.logger.Info("Done waiting for worker")
}

func (group *Group) waitForSignal(workerCancel, signalCancel context.CancelFunc) {
	sig := <-group.signals

	if group.close() {
		return
	}

	group.logger.WithField("signal", sig).Info("Received an os signal")
	group.cleanup(workerCancel, signalCancel)
}

func (group *Group) waitForWorkers(ctx context.Context, workerCancel, signalCancel context.CancelFunc) {
	<-ctx.Done()

	if group.close() {
		return
	}

	group.logger.Info("Received a stop signal from a worker")
	group.cleanup(workerCancel, signalCancel)
}

func (group *Group) close() bool {
	group.mutex.Lock()

	defer group.mutex.Unlock()

	if group.closed {
		return true
	}

	group.closed = true

	return false
}

func (group *Group) cleanup(workerCancel, signalCancel context.CancelFunc) {
	signal.Stop(group.signals)
	close(group.signals)
	signalCancel()
	workerCancel()
}
