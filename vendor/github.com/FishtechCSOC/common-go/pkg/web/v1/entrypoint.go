package web

import (
	"context"
)

type Entrypoint interface {
	Start()
	Shutdown(context.Context)
}
