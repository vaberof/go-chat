package message

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"time"
)

type Message struct {
	Id        domain.MessageId
	SenderId  domain.UserId
	RoomId    domain.RoomId
	Text      string
	CreatedAt time.Time
	EditedAt  time.Time
}
