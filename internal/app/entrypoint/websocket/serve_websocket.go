package websocket

import (
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/middleware"
	"github.com/vaberof/go-chat/internal/domain/auth"
	"github.com/vaberof/go-chat/internal/websocket"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"net/http"
)

func ServeWebsocketRoute(hub *websocket.Hub, authService auth.AuthService, logs *logs.Logs) func(router chi.Router) {
	return func(router chi.Router) {
		router.Handle("/ws", ServeWebsocketHandler(hub, authService, logs))
	}
}

func ServeWebsocketHandler(hub *websocket.Hub, authService auth.AuthService, logs *logs.Logs) http.Handler {
	handler := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		websocket.ServeWebsocket(hub, responseWriter, request)
	})
	return middleware.AuthMiddleware(handler, authService, logs)
}
