package poll

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

// The Poller interface is meant to return data from a source given a time range.
type Poller interface {
	// Type returns the type of Poller
	Type() string
	// Poll accepts a time range to poll with and returns a collection of types.LogEntries if successful.
	Poll(context.Context, TimeRange) []types.LogEntries
}
