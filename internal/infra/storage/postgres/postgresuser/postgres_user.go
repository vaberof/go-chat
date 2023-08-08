package postgresuser

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"time"
)

type User struct {
	Id        domain.UserId
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
