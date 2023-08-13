package http

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
	"time"
)

type createMessageRequestBody struct {
	SenderId domain.UserId `json:"sender_id"`
	RoomId   domain.RoomId `json:"room_id"`
	Text     string        `json:"text"`
}

func (c *createMessageRequestBody) Bind(r *http.Request) error {
	return nil
}

type createMessageResponsePayload struct {
	Id        domain.MessageId `json:"id"`
	SenderId  domain.UserId    `json:"sender_id"`
	RoomId    domain.RoomId    `json:"room_id"`
	Text      string           `json:"text"`
	CreatedAt time.Time        `json:"created_at"`
}

func (h *Handler) CreateMessage() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createMessageReqBody := &createMessageRequestBody{}
		if err := render.Bind(r, createMessageReqBody); err != nil {
			views.RenderJSON(w, r, http.StatusBadRequest, apiv1.Error(InvalidRequestBodyMessage))
			return
		}

		message, err := h.messageService.Create(createMessageReqBody.SenderId, createMessageReqBody.RoomId, createMessageReqBody.Text)
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(createMessageResponsePayload{
			Id:        message.Id,
			SenderId:  message.SenderId,
			RoomId:    message.RoomId,
			Text:      message.Text,
			CreatedAt: message.CreatedAt,
		})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}
