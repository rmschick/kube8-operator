package runnable

import (
	"context"
)

// The Runnable interface is meant just allows a RunnerGroup to create runners and wait for them to complete.
type Runnable interface {
	// Run is meant to handle all necessary lifecycle items for a runnable. This includes startup, runtime, waiting for
	// the given context to be done, and cleanup afterwards.
	Run(context.Context, context.CancelFunc)
	// Wait is just a helper that should block until the Runnable has successfully cleaned up.
	Wait()
}
