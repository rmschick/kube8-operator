package chunk

import (
	"context"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
)

type Chunker interface {
	Chunk(context.Context, <-chan *types.LogEntries, chan<- *types.LogEntries)
}
