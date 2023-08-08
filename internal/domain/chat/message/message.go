package message

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"time"
)

type Message struct {
	Id         uint64
	SenderId   domain.UserId
	Message    string
	ChatRoomId uint64
	CreatedAt  time.Time
}
