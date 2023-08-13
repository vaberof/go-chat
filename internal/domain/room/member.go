package room

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"time"
)

const (
	AdminRole   = "admin"
	RegularRole = "regular"
)

type Member struct {
	UserId   domain.UserId `json:"user_id"`
	RoomId   domain.RoomId `json:"room_id""`
	Nickname string        `json:"nickname"`
	Role     string        `json:"role"`
	JoinedAt time.Time     `json:"joined_at"`
}
