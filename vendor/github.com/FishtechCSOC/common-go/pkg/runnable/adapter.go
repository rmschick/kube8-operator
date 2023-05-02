package runnable

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/common-go/pkg/nursery/v1"
)

var _ Runnable = (*adapterRunnable)(nil)

type adapterRunnable struct {
	job    nursery.Job
	done   chan struct{}
	logger *logrus.Entry
}

func (a *adapterRunnable) Run(ctx context.Context, cancelFunc context.CancelFunc) {
	defer cancelFunc()

	if err := a.job(ctx); err != nil {
		a.logger.WithError(err).Error("encountered error running ")
	}
}

func (a *adapterRunnable) Wait() {
	<-a.done
}

// nolint: ireturn
func NurseryJobToRunnable(job nursery.Job, logger *logrus.Entry) Runnable {
	return &adapterRunnable{
		job:    job,
		done:   make(chan struct{}),
		logger: logger,
	}
}

// nolint: revive
func RunnableToConcurrentJob(runner Runnable) nursery.Job {
	return func(ctx context.Context) error {
		newCtx, cancel := context.WithCancel(ctx)

		runner.Run(newCtx, cancel)
		runner.Wait()

		return newCtx.Err()
	}
}
