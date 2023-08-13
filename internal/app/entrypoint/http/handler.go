package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/middleware"
	websocketroutes "github.com/vaberof/go-chat/internal/app/entrypoint/websocket"
	authservice "github.com/vaberof/go-chat/internal/domain/auth"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/internal/domain/user"
	"github.com/vaberof/go-chat/internal/websocket"
	"github.com/vaberof/go-chat/pkg/logging/logs"
)

type Handler struct {
	hub         *websocket.Hub
	authService authservice.AuthService
	userService user.UserService
	roomService room.RoomService
}

func NewHandler(hub *websocket.Hub, authService authservice.AuthService, userService user.UserService, roomService room.RoomService) *Handler {
	return &Handler{hub: hub, authService: authService, userService: userService, roomService: roomService}
}

func (h *Handler) InitRoutes(router chi.Router, logs *logs.Logs) chi.Router {
	router.Route("/websocket", func(r chi.Router) {
		r.Get("/", websocketroutes.ServeWebsocketHandler(h.hub, h.authService, logs))
	})

	router.Route("/api/v1", func(apiv1 chi.Router) {

		apiv1.Route("/account", func(account chi.Router) {
			account.Post("/register", h.Register())
			account.Post("/login", h.Login())
		})

		apiv1.Route("/users", func(users chi.Router) {
			users.Get("/{user_id}", middleware.AuthMiddleware(h.GetUser(), h.authService, logs))
		})

		apiv1.Route("/rooms", func(rooms chi.Router) {
			rooms.Post("/", middleware.AuthMiddleware(h.CreateRoom(), h.authService, logs))
			rooms.Post("/list", middleware.AuthMiddleware(h.GetUserRooms(), h.authService, logs))
			rooms.Get("/{room_id}/members", middleware.AuthMiddleware(h.GetMembers(), h.authService, logs))
		})

	})

	return router
}
