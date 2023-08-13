package room

import (
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/pkg/domain"
)

type RoomStorage interface {
	Create(creatorId domain.UserId, name, roomType string, members []domain.UserId) (*Room, error)
	Get(roomId domain.RoomId) (*Room, error)
	GetRooms(roomIds []domain.RoomId) ([]*Room, error)
	GetMembers(roomId domain.RoomId) ([]*Member, error)
	GetMessages(roomId domain.RoomId) ([]*message.Message, error)
	List(userId domain.UserId) ([]*Room, error)
	Find(roomId domain.RoomId) error
}
