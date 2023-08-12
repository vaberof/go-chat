package room

import (
	"github.com/vaberof/go-chat/pkg/domain"
)

const (
	PrivateRoomType = "private"
	GroupRoomType   = "group"
)

type Room struct {
	Id        domain.RoomId
	CreatorId domain.UserId
	Name      string
	Type      string
	Members   []domain.UserId
}
