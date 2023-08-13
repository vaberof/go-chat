package room

import (
	"github.com/vaberof/go-chat/pkg/domain"
)

const (
	PrivateRoomType = "private"
	GroupRoomType   = "group"
)

type Room struct {
	Id        domain.RoomId   `json:"id"`
	CreatorId domain.UserId   `json:"creator_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Members   []domain.UserId `json:"members"`
}
