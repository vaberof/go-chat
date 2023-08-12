package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/middleware"
	authservice "github.com/vaberof/go-chat/internal/domain/auth"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/pkg/logging/logs"
)

type Handler struct {
	authService authservice.AuthService
	roomService room.RoomService
}

func NewHandler(authService authservice.AuthService, roomService room.RoomService) *Handler {
	return &Handler{authService: authService, roomService: roomService}
}

func (h *Handler) InitRoutes(router chi.Router, logs *logs.Logs) chi.Router {
	router.Route("/api/v1", func(r chi.Router) {

		r.Route("/account", func(r chi.Router) {
			r.Post("/register", h.Register())
			r.Post("/login", h.Login())
		})

		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", middleware.AuthMiddleware(h.CreateRoom(), h.authService, logs))
		})

	})

	return router
}
