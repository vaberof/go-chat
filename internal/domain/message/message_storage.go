package message

import "github.com/vaberof/go-chat/pkg/domain"

type MessageStorage interface {
	Create(senderId domain.UserId, roomId domain.RoomId, message string) (*Message, error)
}
