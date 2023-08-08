package auth

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"time"
)

type JwtPayload struct {
	UserId    domain.UserId
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func NewPayload(params *TokenParams) *JwtPayload {
	return &JwtPayload{
		UserId:    params.UserId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(params.TokenTtl),
	}
}
