package room

import (
	"github.com/vaberof/go-chat/internal/domain/chat/message"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"github.com/vaberof/go-chat/pkg/domain"
)

type Room struct {
	Id        uint64
	CreatorId domain.UserId
	Name      string
	Members   []*user.User
	Messages  []*message.Message
}
