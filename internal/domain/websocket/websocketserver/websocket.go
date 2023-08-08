package websocketserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/vaberof/go-chat/internal/domain/websocket"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"net/http"
)

type WebsocketServer struct {
	Server  *chi.Mux
	Hub     *websocket.Hub
	config  *WebsocketServerConfig
	logger  *zap.SugaredLogger
	address string
}

func New(config *WebsocketServerConfig, logs *logs.Logs) *WebsocketServer {
	chiServer := chi.NewRouter()
	chiServer.Use(cors.AllowAll().Handler)

	loggerName := "websocket-server"
	logger := logs.WithName(loggerName)

	return &WebsocketServer{
		Server:  chiServer,
		Hub:     websocket.NewHub(logs),
		config:  config,
		logger:  logger,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

func (wsServer *WebsocketServer) StartAsync() *chan error {
	exitChannel := make(chan error)

	go wsServer.Hub.Run()

	wsServer.logger.Infow("Starting WebSocket server")

	go func() {
		err := http.ListenAndServe(wsServer.address, wsServer.Server)
		if err != nil {
			wsServer.logger.Errorw("Failed to start WebSocket server")
			exitChannel <- err
		} else {
			exitChannel <- nil
		}
	}()

	wsServer.logger.Infof("Started WebSocket server at %s", wsServer.address)

	return &exitChannel
}
