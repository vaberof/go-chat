package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vaberof/go-chat/internal/app/entrypoint/http/views"
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
	"time"
)

type getRoomMessagesResponsePayload struct {
	Messages []*getMessageResponse `json:"messages"`
}

type getMessageResponse struct {
	Id        domain.MessageId `json:"id"`
	SenderId  domain.UserId    `json:"sender_id"`
	RoomId    domain.RoomId    `json:"room_id"`
	Text      string           `json:"text"`
	CreatedAt time.Time        `json:"created_at"`
	EditedAt  time.Time        `json:"edited_at"`
}

func (h *Handler) GetMessages() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		messages, err := h.roomService.GetMessages(domain.RoomId(convRoomId))
		if err != nil {
			views.RenderJSON(w, r, http.StatusInternalServerError, apiv1.Error(err.Error()))
			return
		}

		payload, _ := json.Marshal(getRoomMessagesResponsePayload{
			Messages: buildGetMessagesResponse(messages),
		})

		views.RenderJSON(w, r, http.StatusOK, apiv1.Success(payload))
	})
}

func buildGetMessageResponse(message *message.Message) *getMessageResponse {
	return &getMessageResponse{
		Id:        message.Id,
		SenderId:  message.SenderId,
		RoomId:    message.RoomId,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		EditedAt:  message.EditedAt,
	}
}

func buildGetMessagesResponse(messages []*message.Message) []*getMessageResponse {
	messagesResponse := make([]*getMessageResponse, len(messages))

	for i := 0; i < len(messagesResponse); i++ {
		messagesResponse[i] = buildGetMessageResponse(messages[i])
	}

	return messagesResponse
}
