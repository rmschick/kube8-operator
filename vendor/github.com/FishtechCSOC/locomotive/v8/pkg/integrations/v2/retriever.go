package integrations

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

// The Retriever interface is meant to return data from a source as it comes in through a channel.
type Retriever interface {
	// Type returns the type of Retriever
	Type() string
	// Retrieve is meant to handle the business logic for grabbing data and sending it through the given stream.
	// The context is given should indicate to the Retriever when to stop and clean up all resources.
	Retrieve(context.Context, context.CancelFunc, chan<- *types.LogEntries)
	// Wait is just a helper that should block until the Retriever has successfully cleaned up.
	Wait()
}
