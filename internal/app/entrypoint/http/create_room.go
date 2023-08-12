package http

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
)

type createRoomRequestBody struct {
	CreatorId domain.UserId   `json:"creator_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Members   []domain.UserId `json:"members"`
}

type createRoomResponsePayload struct {
	CreatorId domain.UserId   `json:"creator_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Members   []domain.UserId `json:"members"`
}

func (c *createRoomRequestBody) Bind(r *http.Request) error {
	return nil
}

func (h *Handler) CreateRoom() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		createRoomReqBody := &createRoomRequestBody{}

		if err := render.Bind(request, createRoomReqBody); err != nil {
			views.RenderJSON(writer, request, http.StatusBadRequest, apiv1.Error(InvalidRequestBodyMessage))
			return
		}

		room, err := h.roomService.Create(createRoomReqBody.CreatorId, createRoomReqBody.Name, createRoomReqBody.Type, createRoomReqBody.Members)
		if err != nil {
			views.RenderJSON(writer, request, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(createRoomResponsePayload{
			CreatorId: room.CreatorId,
			Name:      room.Name,
			Type:      room.Type,
			Members:   room.Members,
		})

		views.RenderJSON(writer, request, http.StatusOK, apiv1.Success(payload))
	})
}
