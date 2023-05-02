package helpers

import (
	"context"

	"github.com/FishtechCSOC/common-go/pkg/runnable"
	"github.com/arunsworld/nursery"
)

func RunnableToConcurrentJob(runner runnable.Runnable) nursery.ConcurrentJob {
	return func(ctx context.Context, errors chan error) {
		newCtx, cancel := context.WithCancel(ctx)

		runner.Run(newCtx, cancel)
		runner.Wait()

		if err := newCtx.Err(); err != nil {
			errors <- err
		}
	}
}

func FuncToConcurrentJob(f func(ctx context.Context) error) nursery.ConcurrentJob {
	return func(ctx context.Context, errors chan error) {
		if err := f(ctx); err != nil {
			errors <- err
		}
	}
}
