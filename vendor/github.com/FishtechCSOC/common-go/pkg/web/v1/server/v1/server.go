package server

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool"

	"github.com/FishtechCSOC/common-go/pkg/web/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/entrypoints/http/v1"
	"github.com/FishtechCSOC/common-go/pkg/web/v1/logging/v1"
)

// Server is a lifecycle manager for multiple "entrypoints".
type Server struct {
	readiness     *abool.AtomicBool
	done          chan bool
	configuration Configuration
	endpoints     map[string]web.Endpoint
	entrypoints   map[string]web.Entrypoint
	middleware    map[string]web.Middleware
	logger        *logrus.Entry
}

// CreateServer sets up a new server instance.
func CreateServer(endpoints []web.Endpoint, middleware []web.Middleware, configuration Configuration, logger *logrus.Entry) (*Server, error) {
	server := &Server{
		configuration: configuration,
		endpoints:     make(map[string]web.Endpoint),
		entrypoints:   make(map[string]web.Entrypoint),
		middleware:    make(map[string]web.Middleware),
		done:          make(chan bool, 1),
		readiness:     abool.NewBool(false),
	}

	server.logger = logging.SetupServerLogger(logger, "", server)

	gin.SetMode(gin.ReleaseMode)

	endpointNames, err := server.addEndpoints(endpoints...)
	if err != nil {
		return nil, err
	}

	middlewareNames, err := server.addMiddleware(middleware...)
	if err != nil {
		return nil, err
	}

	entrypointNames, err := server.getEntrypoints()
	if err != nil {
		return nil, err
	}

	server.logger.Info("Registered endpoints: " + strings.Join(endpointNames, ","))
	server.logger.Info("Registered middleware: " + strings.Join(middlewareNames, ","))
	server.logger.Info("Registered entrypoints: " + strings.Join(entrypointNames, ","))

	return server, nil
}

// Ready returns an boolean for other consumers to be aware of its status.
func (server *Server) Ready() *abool.AtomicBool {
	return server.readiness
}

// Run is used to start all entrypoints and wait until the context is done before doing cleanup.
func (server *Server) Run(ctx context.Context, cancel context.CancelFunc) {
	defer server.shutdown(cancel)

	server.readiness.Set()
	server.startEntrypoints()
	<-ctx.Done()

	server.drain()
	server.stopEntrypoints(ctx)
}

// Wait is a helper function that does not return until everything was shutdown.
func (server *Server) Wait() {
	server.logger.Info("Began waiting for server to stop")
	<-server.done
	server.logger.Info("Done waiting for server to stop")
}

func (server *Server) startEntrypoints() {
	for _, entrypoint := range server.entrypoints {
		go entrypoint.Start()
	}
}

func (server *Server) stopEntrypoints(ctx context.Context) {
	server.logger.Info("Gracefully shutting down all entrypoints")

	var waitGroup sync.WaitGroup

	for _, entrypoint := range server.entrypoints {
		waitGroup.Add(1)

		go func(ctx context.Context, entrypoint web.Entrypoint, wg *sync.WaitGroup) {
			defer wg.Done()
			entrypoint.Shutdown(ctx)
		}(ctx, entrypoint, &waitGroup)
	}

	waitGroup.Wait()
	server.logger.Info("All entrypoints successfully closed")
}

func (server *Server) drain() {
	server.logger.Info("I have to go...")
	server.readiness.UnSet()

	if timeout := server.configuration.Lifecycle.DrainTimeout; timeout > 0 {
		server.logger.WithField("drainTimeout", timeout).Info("Waiting for incoming requests to cease")
		time.Sleep(timeout)
	}
}

func (server *Server) shutdown(cancel context.CancelFunc) {
	server.logger.Info("Shutting down...")
	close(server.done)
	cancel()
	server.logger.Info("Successfully shutdown")
}

func (server *Server) addMiddleware(middleware ...web.Middleware) ([]string, error) {
	names := make([]string, 0, len(middleware))

	for _, middle := range middleware {
		if _, ok := server.middleware[middle.Name()]; ok {
			return nil, errors.New(middle.Name() + " already registered as a middleware")
		}

		server.middleware[middle.Name()] = middle
		names = append(names, middle.Name())
	}

	return names, nil
}

func (server *Server) addEndpoints(endpoints ...web.Endpoint) ([]string, error) {
	names := make([]string, 0, len(endpoints))

	for _, endpoint := range endpoints {
		if _, ok := server.endpoints[endpoint.Name()]; ok {
			return nil, errors.New(endpoint.Name() + " already registered as a endpoint")
		}

		server.endpoints[endpoint.Name()] = endpoint
		names = append(names, endpoint.Name())
	}

	return names, nil
}

func (server *Server) getEntrypoints() ([]string, error) {
	names := make([]string, 0, len(server.configuration.Entrypoints))

	for name, entrypointConfiguration := range server.configuration.Entrypoints {
		entrypoint, err := server.createEntrypoint(name, entrypointConfiguration)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create entrypoint "+name)
		}

		server.entrypoints[name] = entrypoint

		names = append(names, name)
	}

	return names, nil
}

// nolint: ireturn
func (server *Server) createEntrypoint(name string, configuration EntrypointConfiguration) (web.Entrypoint, error) {
	switch {
	case configuration.HTTP != nil:
		return http.CreateEntrypoint(name, server.endpoints, server.middleware, *configuration.HTTP, server.configuration.Lifecycle.GracefulTimeout, server.logger), nil
	default:
		return nil, errors.New("unable to create entrypoint due to empty configuration: " + name)
	}
}
