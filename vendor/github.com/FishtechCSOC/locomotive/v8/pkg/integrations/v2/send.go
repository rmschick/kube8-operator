package integrations

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

// SendToStream is meant to be used by a Retriever to respect a closed context in the case where it is unable to send
// data through the types.LogEntries stream because the buffer is full (or does not exist) and helping avoid orphaned
// goroutines.
func SendToStream(ctx context.Context, stream chan<- *types.LogEntries, entries *types.LogEntries) error {
	select {
	case stream <- entries:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
