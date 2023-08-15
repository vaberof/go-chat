package websocket

import (
	"encoding/json"
	"github.com/vaberof/go-chat/pkg/domain"
)

type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type MessagePayload struct {
	SenderId domain.UserId `json:"sender_id"`
	RoomId   domain.RoomId `json:"room_id"`
	Text     string        `json:"text"`
}

type JoinRoomPayload struct {
	RoomId domain.RoomId `json:"room_id"`
}

func (m *Message) Encode() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return bytes
}

func (m *MessagePayload) Encode() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return bytes
}
