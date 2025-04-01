package instrumentation

import (
	"go.opencensus.io/tag"
)

// nolint: gochecknoglobals
var (
	PathTag   = tag.MustNewKey("path")
	HostTag   = tag.MustNewKey("host")
	MethodTag = tag.MustNewKey("method")
)
