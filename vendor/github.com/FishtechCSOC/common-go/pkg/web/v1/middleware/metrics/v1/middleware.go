package metrics

import (
	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/FishtechCSOC/common-go/pkg/metrics/instrumentation"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/logging/v1"
)

const (
	namespace   = "metrics"
	packageName = "github.com/FishtechCSOC/common-go/pkg/server/middleware/metrics"
)

type Middleware struct{}

func CreateMiddleware() *Middleware {
	return &Middleware{}
}

func (middleware *Middleware) Name() string {
	return namespace
}

func (middleware *Middleware) Handle(ctx *gin.Context) {
	logger := logging.SetupMiddlewareLogger(namespace, packageName, ctx)

	tags := []tag.Mutator{
		tag.Upsert(instrumentation.PathTag, ctx.Request.URL.Path),
		tag.Upsert(instrumentation.HostTag, ctx.Request.Host),
		tag.Upsert(instrumentation.MethodTag, ctx.Request.Method),
	}

	err := middleware.recordRequestMetrics(ctx, tags)
	if err != nil {
		logger.WithError(err).Info("failed to record request metrics")
	}

	ctx.Next()

	err = middleware.recordResponseMetrics(ctx, tags)
	if err != nil {
		logger.WithError(err).Info("failed to record response metrics")
	}
}

func (middleware *Middleware) recordRequestMetrics(ctx *gin.Context, tags []tag.Mutator) error {
	measurements := []stats.Measurement{
		ochttp.ServerRequestCount.M(1),
		ochttp.ServerRequestBytes.M(ctx.Request.ContentLength),
	}

	err := stats.RecordWithTags(ctx, tags, measurements...)
	if err != nil {
		return err
	}

	return nil
}

func (middleware *Middleware) recordResponseMetrics(ctx *gin.Context, tags []tag.Mutator) error {
	measurements := []stats.Measurement{
		ochttp.ServerResponseBytes.M(int64(ctx.Writer.Size())),
	}

	err := stats.RecordWithTags(ctx, tags, measurements...)
	if err != nil {
		return err
	}

	return nil
}
