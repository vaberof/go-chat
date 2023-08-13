package user

import (
	"github.com/vaberof/go-chat/pkg/domain"
)

type User struct {
	Id       domain.UserId
	Username string
	Password string
	Rooms    []domain.RoomId
	Messages []domain.MessageId
}
