package metrics

import (
	"go.opencensus.io/tag"
)

// The following tags are applied to stats recorded by this package.
// nolint: gochecknoglobals
var (
	// CDP Customer Name.
	CustomerNameTag = tag.MustNewKey("customer_name")

	// CDP Customer ID.
	CustomerIDTag = tag.MustNewKey("customer_id")

	// CDP Log Type.
	LogTypeTag = tag.MustNewKey("log_type")

	// CDP Data Type.
	DataTypeTag = tag.MustNewKey("data_type")

	// Success/Fail Status.
	StatusTag = tag.MustNewKey("status")
)
