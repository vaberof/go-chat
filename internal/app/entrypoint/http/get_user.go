package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

type getUserResponsePayload struct {
	Id       domain.UserId   `json:"id"`
	Username string          `json:"username"`
	Rooms    []domain.RoomId `json:"rooms"`
}

func (h *Handler) GetUser() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "user_id")
		if userId == "" {
			views.RenderJSON(w, r, http.StatusNotFound, apiv1.Error("Resource not found"))
			return
		}

		convUserId, err := strconv.Atoi(userId)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(fmt.Sprintf("Failed to convert userId: %v", err)))
			return
		}

		user, err := h.userService.Get(domain.UserId(convUserId))
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(getUserResponsePayload{
			Id:       user.Id,
			Username: user.Username,
			Rooms:    user.Rooms,
		})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
