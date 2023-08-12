package websocket

import (
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/middleware"
	"github.com/vaberof/go-chat/internal/domain/auth"
	"github.com/vaberof/go-chat/internal/websocket"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"net/http"
)

func ServeWebsocketHandler(hub *websocket.Hub, authService auth.AuthService, logs *logs.Logs) http.HandlerFunc {
	handler := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		websocket.ServeWebsocket(hub, responseWriter, request)
	})

	return middleware.AuthMiddleware(handler, authService, logs)
}
