package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

type getMembersResponsePayload struct {
	Members []*room.Member `json:"members"`
}

func (h *Handler) GetMembers() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("url: %v\n", r.URL.Path)
		roomId := chi.URLParam(r, "room_id")
		if roomId == "" {
			views.RenderJSON(w, r, http.StatusNotFound, apiv1.Error("Resource not found"))
			return
		}

		convRoomId, err := strconv.Atoi(roomId)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(fmt.Sprintf("Failed to convert roomId: %v", err)))
			return
		}

		members, err := h.roomService.GetMembers(domain.RoomId(convRoomId))
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(getMembersResponsePayload{
			Members: members,
		})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
