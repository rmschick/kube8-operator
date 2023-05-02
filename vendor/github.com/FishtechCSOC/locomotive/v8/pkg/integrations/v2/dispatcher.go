package integrations

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

// Dispatcher interface is meant to take data and pass it to another datastore.
type Dispatcher interface {
	// Type returns the type of Dispatcher
	Type() string
	// Dispatch should massage and push the given types.LogEntries to a target. The context is given should indicate to
	// the Dispatcher when to stop and give up on the dispatches it is processing.
	Dispatch(context.Context, *types.LogEntries) types.Status
}
