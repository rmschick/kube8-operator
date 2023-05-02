package metrics

import (
	"github.com/sirupsen/logrus"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2"
)

func SetupCursor(cursor integrations.Cursor, logger *logrus.Entry, tags map[string]string) *Cursor {
	tagKeys := make([]tag.Key, 0, len(tags))
	for k := range tags {
		tagKeys = append(tagKeys, tag.MustNewKey(k))
	}

	err := RegisterView(tagKeys...)
	if err != nil {
		panic(err)
	}

	return CreateCursor(cursor, logger, tags)
}
