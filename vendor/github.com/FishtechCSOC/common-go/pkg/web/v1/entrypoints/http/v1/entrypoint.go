package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/common-go/pkg/web/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/logging/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/middleware/accesslog/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/middleware/injector/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/middleware/metrics/v1"
)

type Entrypoint struct {
	name            string
	configuration   Configuration
	gracefulTimeout time.Duration
	httpServer      *http.Server
	router          *gin.Engine
	logger          *logrus.Entry
}

func CreateEntrypoint(name string, endpoints map[string]web.Endpoint, middleware map[string]web.Middleware, configuration Configuration, gracefulTimeout time.Duration, logger *logrus.Entry) *Entrypoint {
	entrypoint := &Entrypoint{
		name:            name,
		configuration:   configuration,
		gracefulTimeout: gracefulTimeout,
		httpServer: &http.Server{
			Addr:              configuration.Host + ":" + strconv.Itoa(configuration.Port),
			ReadTimeout:       configuration.ReadTimeout,
			ReadHeaderTimeout: configuration.ReadHeaderTimeout,
			WriteTimeout:      configuration.WriteTimeout,
			IdleTimeout:       configuration.IdleTimeout,
		},
		router: gin.New(),
	}

	entrypoint.logger = logging.SetupEntrypointLogger(logger, name, entrypoint)

	entrypoint.addDefaultMiddleware()

	middlewareNames := entrypoint.addMiddleware(middleware)
	entrypoint.logger.Info("Registered middleware: " + strings.Join(middlewareNames, ","))

	entrypoint.router.Use(gin.Recovery())

	endpointsNames := entrypoint.addEndpoints(endpoints)
	entrypoint.logger.Info("Registered endpoints: " + strings.Join(endpointsNames, ","))

	entrypoint.httpServer.Handler = entrypoint.router

	return entrypoint
}

func (entrypoint *Entrypoint) Start() {
	entrypoint.logger.WithField("address", entrypoint.httpServer.Addr).Info("Starting server")

	err := entrypoint.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		entrypoint.logger.WithError(err).Error("Error creating server")
	}
}

func (entrypoint *Entrypoint) Shutdown(ctx context.Context) {
	timeout := entrypoint.gracefulTimeout
	cancelCtx, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	entrypoint.gracefulShutdown(cancelCtx)
	entrypoint.logger.Info("Entrypoint is shutdown...")
}

func (entrypoint *Entrypoint) gracefulShutdown(ctx context.Context) {
	if err := entrypoint.httpServer.Shutdown(ctx); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			entrypoint.logger.WithError(err).Error("Wait server shutdown is overdue")
		default:
			entrypoint.logger.WithError(err).Error("Failed to shutdown http server due to unrecognized error")
		}

		err = entrypoint.httpServer.Close()
		if err != nil {
			entrypoint.logger.WithError(err).Error("Closing http server failed due to error")
		}
	}
}

func (entrypoint *Entrypoint) addMiddleware(middleware map[string]web.Middleware) []string {
	names := make([]string, 0, len(middleware))

	for _, middlewareName := range entrypoint.configuration.Middleware {
		middle, ok := middleware[middlewareName]

		if !ok {
			entrypoint.logger.WithField("middleware", middlewareName).Error("no middleware with given name was registered with server, ignoring")

			continue
		}

		if contains(names, middlewareName) {
			entrypoint.logger.WithField("middleware", middlewareName).Info("Duplicate middleware configured, ignoring")

			continue
		}

		names = append(names, middlewareName)

		entrypoint.router.Use(middle.Handle)
	}

	return names
}

func (entrypoint *Entrypoint) addEndpoints(endpoints map[string]web.Endpoint) []string {
	names := make([]string, 0, len(endpoints))

	for _, endpointName := range entrypoint.configuration.Endpoints {
		endpoint, ok := endpoints[endpointName]

		if !ok {
			entrypoint.logger.WithField("endpoint", endpointName).Error("no endpoint with given name was registered with server, ignoring")

			continue
		}

		if contains(names, endpointName) {
			entrypoint.logger.WithField("endpoint", endpointName).Info("Duplicate endpoint configured, ignoring")

			continue
		}

		names = append(names, endpointName)

		endpoint.AddRoutes(entrypoint.router.Group("/" + endpoint.Name()))
	}

	return names
}

func (entrypoint *Entrypoint) addDefaultMiddleware() {
	injectorMiddle := injector.CreateMiddleware(injector.Configuration{
		Logger: entrypoint.logger,
		Router: entrypoint.router,
	})
	entrypoint.router.Use(injectorMiddle.Handle)

	accesslogMiddle := accesslog.CreateMiddleware(entrypoint.configuration.Accesslog)
	entrypoint.router.Use(accesslogMiddle.Handle)

	if entrypoint.configuration.Metrics.Enabled {
		metricsMiddle := metrics.CreateMiddleware()
		entrypoint.router.Use(metricsMiddle.Handle)
	}

	entrypoint.router.Use(gin.Recovery())
}

func contains(source []string, target string) bool {
	for _, value := range source {
		if value == target {
			return true
		}
	}

	return false
}
