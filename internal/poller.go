package internal

import (
	"context"

	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/poll/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
	"github.com/sirupsen/logrus"
)

var _ poll.Poller = (*Poller)(nil)

type Poller struct {
	logger   *logrus.Entry
	metadata types.Metadata
}

func CreatePoller(metadata types.Metadata, logger *logrus.Entry) *Poller {
	poller := &Poller{
		metadata: metadata,
	}

	poller.logger = poll.SetupPollerLogger(string(metadata.DataType), poller, logger)

	return poller
}

func (poller *Poller) Poll(_ context.Context, _ poll.TimeRange) []types.LogEntries {
	// Add your API calls here
	return []types.LogEntries{
		{
			Metadata: poller.metadata,
			Entries: []types.LogEntry{
				{
					Log: "",
				},
			},
		},
		{
			Metadata: poller.metadata,
			Entries: []types.LogEntry{
				{
					Log: "",
				},
			},
		},
	}
}

func (poller *Poller) Type() string {
	return build.Program
}
