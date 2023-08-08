package auth

import (
	"errors"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"github.com/vaberof/go-chat/pkg/auth"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/xpassword"
	"time"
)

const defaultTokenTtl = 7 * 24 * time.Hour

type AuthService interface {
	Register(username, password string) (*user.User, error)
	Login(username, password string) (*auth.JwtToken, error)
	VerifyToken(token string) (*domain.UserId, error)
}

type authServiceImpl struct {
	userService user.UserService
	config      *AuthServiceConfig
}

func NewAuthService(userService user.UserService, config *AuthServiceConfig) AuthService {
	return &authServiceImpl{userService: userService, config: config}
}

func (service *authServiceImpl) Register(username, password string) (*user.User, error) {
	err := service.userService.FindByUsername(username)
	if err == nil {
		return nil, errors.New("user with given username already exists")
	}

	user, err := service.userService.Create(username, password)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (service *authServiceImpl) Login(username, password string) (*auth.JwtToken, error) {
	user, err := service.userService.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	err = xpassword.Check(password, user.Password)
	if err != nil {
		return nil, errors.New("Incorrect password")
	}

	tokenParams := &auth.TokenParams{
		UserId:   user.Id,
		TokenTtl: defaultTokenTtl,
	}

	token, err := auth.CreateToken(tokenParams, auth.SecretKey(service.config.SecretKey))
	if err != nil {
		return nil, err
	}

	jwtToken := auth.JwtToken(token)

	return &jwtToken, nil
}

func (service *authServiceImpl) VerifyToken(token string) (*domain.UserId, error) {
	payload, err := auth.VerifyToken(token, auth.SecretKey(service.config.SecretKey))
	if err != nil {
		return nil, err
	}

	user, err := service.userService.Get(payload.UserId)
	if err != nil {
		return nil, err
	}

	return &user.Id, nil
}
