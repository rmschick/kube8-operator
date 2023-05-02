package firestore

import (
	"github.com/FishtechCSOC/locomotive/v8/pkg/helpers"
)

const (
	CursorType = "firestore"

	NoCursorFoundError        helpers.Error = "found no matching cursors"
	MultipleCursorFoundError  helpers.Error = "found multiple matching cursors"
	CreateCursorDocumentError helpers.Error = "failed to create cursor document"

	v1CollectionRef = "cdp.v1.cursor"
)
