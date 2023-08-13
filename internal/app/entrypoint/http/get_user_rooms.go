package http

import (
	"encoding/json"
	"fmt"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

type getUserRoomsResponsePayload struct {
	Rooms []*room.Room `json:"rooms"`
}

func (h *Handler) GetUserRooms() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		if userId == "" {
			views.RenderJSON(w, r, http.StatusBadRequest, apiv1.Error("Missing required query parameter 'userId'"))
			return
		}

		convUserId, err := strconv.Atoi(userId)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(fmt.Sprintf("Failed to convert userId: %v", err)))
			return
		}

		rooms, err := h.roomService.List(domain.UserId(convUserId))
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(getUserRoomsResponsePayload{
			Rooms: rooms,
		})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
