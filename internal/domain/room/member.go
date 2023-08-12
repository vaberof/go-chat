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
	UserId   domain.UserId
	RoomId   domain.RoomId
	Nickname string
	Role     string
	JoinedAt time.Time
}
