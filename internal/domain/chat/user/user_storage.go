package user

import "github.com/vaberof/go-chat/pkg/domain"

type UserStorage interface {
	Create(username, password string) (*User, error)
	Get(userIf domain.UserId) (*User, error)
	GetByUsername(username string) (*User, error)
	FindByUsername(username string) error
}
