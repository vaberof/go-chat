package httpserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/vaberof/go-chat/pkg/http/httpserver/middleware/logging"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"net/http"
)

type HttpServer struct {
	Server  *chi.Mux
	config  *HttpServerConfig
	logger  *zap.SugaredLogger
	address string
}

func New(config *HttpServerConfig, logs *logs.Logs) *HttpServer {
	loggingMw := logging.New(logs)

	chiServer := chi.NewRouter()
	chiServer.Use(
		loggingMw.Handler,
		cors.AllowAll().Handler)

	return &HttpServer{
		Server:  chiServer,
		config:  config,
		logger:  loggingMw.Logger,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

func (server *HttpServer) StartAsync() *chan error {
	exitChannel := make(chan error)
	server.logger.Infow("Starting http server")
	go func() {
		err := http.ListenAndServe(server.address, server.Server)
		if err != nil {
			server.logger.Errorw("Failed to start HTTP server")
			exitChannel <- err
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Infof("Started HTTP server at %s", server.address)
	return &exitChannel
}

func (server *HttpServer) GetLogger() *zap.SugaredLogger {
	return server.logger
}
