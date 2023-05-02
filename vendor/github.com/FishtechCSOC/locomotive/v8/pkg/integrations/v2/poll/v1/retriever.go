package poll

import (
	"context"
	"sync"
	"time"

	cyderesErrors "github.com/FishtechCSOC/common-go/pkg/errors/v1"
	"github.com/TrevinTeacutter/goback"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

const (
	cursorFormat = time.RFC3339Nano
)

const WindowTooSmallError = cyderesErrors.BasicError("calculated poll time window was smaller than configured initial window")

var _ integrations.Retriever = (*Retriever)(nil)

// Retriever handles the management of the actual polling lifecycle.
type Retriever struct {
	configuration Configuration
	poller        Poller
	cursor        integrations.Cursor
	polling       *time.Ticker
	backoff       *goback.JitterBackoff
	done          chan bool
	logger        *logrus.Entry
}

// CreateRetriever sets up a new Retriever instance.
func CreateRetriever(configuration Configuration, poller Poller, cursor integrations.Cursor, logger *logrus.Entry) (*Retriever, error) {
	retriever := &Retriever{
		configuration: configuration,
		poller:        poller,
		cursor:        cursor,
		done:          make(chan bool, 1),
	}

	retriever.logger = integrations.SetupRetrieverLogger(poller.Type(), retriever, logger)

	if configuration.Lifecycle.PollingInterval > 0 {
		retriever.polling = time.NewTicker(configuration.Lifecycle.PollingInterval)
	}

	if retriever.configuration.Lifecycle.Backoff.Minimum >= 0 ||
		retriever.configuration.Lifecycle.Backoff.Maximum >= 0 ||
		retriever.configuration.Lifecycle.Backoff.Minimum >= retriever.configuration.Lifecycle.Backoff.Maximum ||
		retriever.configuration.Lifecycle.Backoff.Factor > 0.0 {
		retriever.backoff = &goback.JitterBackoff{
			SimpleBackoff: goback.SimpleBackoff{
				Min:    retriever.configuration.Lifecycle.Backoff.Minimum,
				Max:    retriever.configuration.Lifecycle.Backoff.Maximum,
				Factor: retriever.configuration.Lifecycle.Backoff.Factor,
			},
		}
	}

	return retriever, nil
}

func (retriever *Retriever) Type() string {
	return retriever.poller.Type()
}

// Retrieve handles the main lifecycle loop using a done context as the closing case.
func (retriever *Retriever) Retrieve(ctx context.Context, cancel context.CancelFunc, stream chan<- *types.LogEntries) {
	// nolint: contextcheck
	defer retriever.close(cancel)

	retriever.logger.Info("Doing initial poll")
	retriever.process(ctx, stream)

	if retriever.polling == nil {
		retriever.logger.Info("Poller was configured without polling interval, closing")

		return
	}

	retriever.logger.Info("Starting polling")

	for {
		select {
		case <-ctx.Done():
			retriever.logger.WithField("reason", ctx.Err()).Info("context was closed, closing")

			return
		case pollTime := <-retriever.polling.C:
			retriever.logger.WithField("time", pollTime).Debug("Tick")
			retriever.process(ctx, stream)
		}
	}
}

// Wait blocks until the poller and dispatcher are successfully cleaned up.
func (retriever *Retriever) Wait() {
	retriever.logger.Info("Began waiting for polling ticker to stop")
	<-retriever.done
	retriever.logger.Info("Done waiting for polling ticker to stop")
}

func (retriever *Retriever) process(ctx context.Context, stream chan<- *types.LogEntries) {
	timeRange, err := retriever.calculateWindow(ctx, time.Now().UTC().Add(-retriever.configuration.Lifecycle.WindowOffset))
	if err != nil {
		switch errors.Is(err, WindowTooSmallError) {
		case false:
			retriever.logger.WithError(err).Error("Failed to calculate the poll time window")
		default:
			retriever.logger.Debug(WindowTooSmallError.Error())
		}

		return
	}

	retriever.logger.WithFields(logrus.Fields{
		"window.start": timeRange.Start,
		"window.end":   timeRange.End,
	}).Info("Calculated Window, Polling...")

	if !retriever.processEntries(ctx, stream, timeRange) {
		retriever.logger.Info("Failed to poll, leaving cursor where it is and backing off")

		if err = goback.Wait(retriever.backoff); err != nil {
			retriever.logger.WithError(err).Error("failed to backoff")
		}

		return
	}

	retriever.backoff.Reset()

	if err = retriever.cursor.Store(ctx, timeRange.End.Format(cursorFormat)); err != nil {
		retriever.logger.WithError(err).Error("Failed to store cursor")

		return
	}
}

func (retriever *Retriever) processEntries(ctx context.Context, stream chan<- *types.LogEntries, timeRange TimeRange) bool {
	allEntries := retriever.poller.Poll(ctx, timeRange)

	if allEntries == nil {
		retriever.logger.Info("Nil entries were returned by poller")

		return false
	}

	if len(allEntries) <= 0 {
		retriever.logger.Info("No entries to push, moving timestamp")

		return true
	}

	group := &sync.WaitGroup{}
	resultStream := make(chan types.Status, len(allEntries))

	for _, entries := range allEntries {
		// Put uuid onto batchID so that we can merge the trace if it doesn't exist
		if entries.Metadata.BatchID == "" {
			entries.Metadata.BatchID = uuid.New().String()
		}

		once := &sync.Once{}
		ref := entries
		ref.Status = retriever.statusFunction(group, once, resultStream)

		group.Add(1)

		if err := integrations.SendToStream(ctx, stream, &ref); err != nil {
			retriever.logger.WithError(err).Info("Failed to push data to stream")

			return false
		}
	}

	group.Wait()
	close(resultStream)

	for status := range resultStream {
		if status == types.Failure {
			return false
		}
	}

	return true
}

func (retriever *Retriever) statusFunction(group *sync.WaitGroup, once *sync.Once, stream chan<- types.Status) func(status types.Status) {
	return func(status types.Status) {
		once.Do(func() {
			defer group.Done()

			stream <- status
		})
	}
}

func (retriever *Retriever) calculateSlidingPoint(difference time.Duration) time.Duration {
	if retriever.configuration.Lifecycle.SlidingWindowSize < difference {
		return difference
	}

	return retriever.configuration.Lifecycle.SlidingWindowSize
}

func (retriever *Retriever) calculateWindow(ctx context.Context, currentTime time.Time) (TimeRange, error) {
	var cursorTime time.Time

	timestamp, err := retriever.cursor.Load(ctx)
	if err != nil {
		return TimeRange{}, errors.Wrap(err, "Failed to load the poll cursor")
	}

	if timestamp == "" {
		return TimeRange{
			Start: currentTime.Add(-retriever.calculateSlidingPoint(retriever.configuration.Lifecycle.InitialWindow)),
			End:   currentTime,
		}, nil
	}

	cursorTime, err = time.Parse(cursorFormat, timestamp)
	if err != nil {
		return TimeRange{}, errors.Wrap(err, "Failed to parse the poll cursor")
	}

	difference := currentTime.Sub(cursorTime)
	if difference < retriever.configuration.Lifecycle.InitialWindow {
		return TimeRange{}, WindowTooSmallError
	}

	timeRange := TimeRange{
		Start: currentTime.Add(-retriever.calculateSlidingPoint(currentTime.Sub(cursorTime))),
		End:   currentTime,
	}

	if retriever.configuration.Lifecycle.MaxWindowSize > 0 && difference > retriever.configuration.Lifecycle.MaxWindowSize {
		timeRange.End = cursorTime.Add(retriever.configuration.Lifecycle.MaxWindowSize)
	}

	return timeRange, nil
}

func (retriever *Retriever) close(cancel context.CancelFunc) {
	retriever.logger.Info("Stopping tickers and closing channels")

	if retriever.polling != nil {
		retriever.polling.Stop()
	}

	ctx, timeoutCancel := context.WithTimeout(context.Background(), retriever.configuration.Lifecycle.ShutdownTimeout)

	defer timeoutCancel()
	retriever.cursor.Close(ctx)
	close(retriever.done)
	cancel()
	retriever.logger.Info("Tickers stopped and channels closed")
}
