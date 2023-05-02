package main

import (
	"context"
	"runtime"
	"strconv"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/FishtechCSOC/common-go/pkg/build"
	"github.com/FishtechCSOC/common-go/pkg/logging/v1"
	"github.com/FishtechCSOC/common-go/pkg/runnable"
	"github.com/FishtechCSOC/locomotive/v8/pkg/chunker/linear/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/cursors/firestore/v2"
	"github.com/FishtechCSOC/locomotive/v8/pkg/cursors/metrics/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/cursors/throttled/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/dispatchers/cdp/ingestion/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/poll/v1"
	"github.com/FishtechCSOC/locomotive/v8/pkg/integrations/v2/runner"
	"github.com/FishtechCSOC/locomotive/v8/pkg/setup"
	"github.com/spf13/viper"

	"github.com/FishtechCSOC/locomotive-collector-template/internal"
)

func main() {
	var config internal.Configuration

	setup.BuildConfiguration(
		viper.New(),
		&config,
		*internal.ConfigurationDefaults,
		"metadata.customer.id",
		"metadata.tenant.id",
	)

	ctx := context.Background()
	logger := logging.CreateEntry(logging.CreateLogger(config.Logging), config.Logging)

	defer logger.WithField("goroutines", strconv.Itoa(runtime.NumGoroutine())).Info("Exiting")

	dispatcher := ingestion.SetupBucketHandle(ctx, config.Dispatcher, config.Metadata, logger)
	chunker := linear.CreateChunker(config.Chunk, logger)
	pollerShard := "poller"
	pollerCursor := metrics.SetupCursor(firestore.SetupCursor(ctx, config.Metadata, config.Cursor.Firestore, pollerShard, logger), logger, map[string]string{
		"collector":          build.Program,
		"collector_instance": config.Metadata.Instance,
		"customer_id":        config.Metadata.Customer.ID,
		"customer_name":      config.Metadata.Customer.Name,
		"shard":              pollerShard,
	})

	streamerShard := "streamer"

	childCursor, err := throttled.NewCursor(config.Cursor.Throttled, firestore.SetupCursor(ctx, config.Metadata, config.Cursor.Firestore, streamerShard, logger))
	if err != nil {
		panic(err)
	}

	streamerCursor := metrics.SetupCursor(childCursor, logger, map[string]string{
		"collector":          build.Program,
		"collector_instance": config.Metadata.Instance,
		"customer_id":        config.Metadata.Customer.ID,
		"customer_name":      config.Metadata.Customer.Name,
		"shard":              streamerShard,
	})

	poller, err := poll.CreateRetriever(config.Poller, internal.CreatePoller(config.Metadata, logger), pollerCursor, logger)
	if err != nil {
		panic(err)
	}

	streamer := internal.CreateStreamer(config.Metadata, config.Streamer, streamerCursor, logger)
	pollRunner := runner.CreateRunner("Poll Collector", config.Runner, poller, dispatcher, chunker, logger)
	streamRunner := runner.CreateRunner("Stream Collector", config.Runner, streamer, dispatcher, chunker, logger)

	prometheusExporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		panic(err)
	}

	setup.RegisterMetricViews(logger)

	httpServer := setup.BuildHTTPServer(config.Server, prometheusExporter, logger)
	groupRunner := runnable.CreateGroup(logger, httpServer, pollRunner, streamRunner)

	groupRunner.Run(ctx)
	groupRunner.Wait()
}
