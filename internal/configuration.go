package internal

import (
	"time"

	"github.com/FishtechCSOC/common-go/pkg/configuration/v1"
	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/server/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/cursors/firestore/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/cursors/throttled/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/dispatchers/cdp/ingestion/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/poll/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/runner"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/runner/chunk"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/stream"
	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
	"github.com/alecthomas/units"
)

type CursorConfiguration struct {
	Firestore firestore.Configuration `mapstructure:"firestore"`
	Throttled throttled.Configuration `mapstructure:"throttled"`
}

// Configuration is the amalgamation of various configurations that may be needed.
type Configuration struct {
	Metadata types.Metadata        `mapstructure:"metadata"`
	Logging  logging.Configuration `mapstructure:"logging"`
	// Regardless of whether you use Streamer/Poller, this should be renamed to `Retriever`
	// as use a lowercase version for tags because it is what the library uses consistently
	Streamer   stream.Configuration    `mapstructure:"streamer"`
	Poller     poll.Configuration      `mapstructure:"poller"`
	Dispatcher ingestion.Configuration `mapstructure:"dispatcher"`
	Chunk      chunk.Configuration     `mapstructure:"chunk"`
	Runner     runner.Configuration    `mapstructure:"runner"`
	Server     server.Configuration    `mapstructure:"server"`
	Cursor     CursorConfiguration     `mapstructure:"cursor"`
}

// nolint: gochecknoglobals, gomnd
var (
	SourceDefault = configuration.CreateDefaults("source").WithFields(map[string]any{
		"agent": types.CreateUserAgent(),
		"type":  "API Hybrid",
	})

	MetadataDefaults = configuration.CreateDefaults("metadata").WithFields(map[string]any{
		"labels":       map[string]string{},
		"dataType":     types.DataType("YOUR_TYPE_GOES_HERE"),
		"destinations": []types.Destination{types.Chronicle},
	}).WithChildren(SourceDefault)

	ChunkDefaults = chunk.ConfigurationDefaults.WithFields(map[string]any{
		"sizeInBytes": units.Megabyte * 10,
	})

	CursorDefaults = configuration.CreateDefaults("cursor").WithChildren(
		firestore.ConfigurationDefaults.WithPrefix("firestore"),
		configuration.CreateDefaults("throttled").WithFields(map[string]any{
			"throttleOffset": time.Second * 15,
		}),
	)

	PollerDefaults = poll.ConfigurationDefaults.WithPrefix("poller")

	StreamerDefaults = stream.ConfigurationDefaults.WithPrefix("streamer")

	ConfigurationDefaults = configuration.CreateDefaults("").WithChildren(
		// How to handle defaults on the pieces that need to be flexed?
		MetadataDefaults,
		logging.ConfigurationDefaults,
		StreamerDefaults,
		PollerDefaults,
		ingestion.ConfigurationDefaults,
		ChunkDefaults,
		runner.ConfigurationDefaults,
		server.ConfigurationDefaults,
		CursorDefaults,
	)
)
