package integrations

import (
	"context"
)

type Cursor interface {
	// Store stores the given token as the cursor and return an error in the event it is unable.
	Store(context.Context, string) error
	// Load returns the token stored in this Cursor and returns an error if it is unable to retrieve it.
	Load(context.Context) (string, error)
	// Test returns nil or an error depending on whether the Cursor backend is reachable or not.
	Test(context.Context) error
	// Close Cursor by syncing any cached values and clean up resources.
	Close(context.Context)
}
