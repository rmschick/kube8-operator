// Package nursery implements "structured concurrency" in Go.
//
// It's based on this blog post: https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
package nursery

import (
	"context"
)

// Job is a blocking function that is meant to respect closure of the given context and return an error in case the
// function encounters an unrecoverable error.
type Job func(context.Context) error

// Nursery is expected to be a stateful struct that helps manage running Jobs concurrently in a structured manner.
type Nursery interface {
	// Add is expected to take any number of jobs and run them in a goroutine.
	Add(jobs ...Job)
	// Wait is expected to block until the nursery is closed through some condition but will continue to block until
	// all managed goroutines successfully exit.
	Wait() error
	// Active should just return how many goroutines are managed by this nursery currently for the purpose of
	// visibility.
	Active() int
	// Close is an escape valve that should trigger a close condition within the nursery.
	Close()
}
