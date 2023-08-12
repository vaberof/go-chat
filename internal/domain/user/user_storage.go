package user

import (
	"github.com/vaberof/go-chat/pkg/domain"
)

type UserStorage interface {
	Create(username, password string) (*User, error)
	Get(userId domain.UserId) (*User, error)
	GetUsers(userIds []domain.UserId) ([]*User, error)
	GetByUsername(username string) (*User, error)
	FindByUsername(username string) error
}
