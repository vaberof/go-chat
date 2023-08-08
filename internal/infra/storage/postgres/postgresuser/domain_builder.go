package postgresuser

import "github.com/vaberof/go-chat/internal/domain/chat/user"

func buildDomainUser(postgresUser *User) *user.User {
	return &user.User{
		Id:       postgresUser.Id,
		Username: postgresUser.Username,
		Password: postgresUser.Password,
	}
}
