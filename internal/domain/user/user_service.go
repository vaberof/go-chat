package user

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"github.com/vaberof/go-chat/pkg/xpassword"
	"go.uber.org/zap"
)

type UserService interface {
	Create(username, password string) (*User, error)
	Get(userId domain.UserId) (*User, error)
	GetByUsername(username string) (*User, error)
	FindByUsername(username string) error
}

type userServiceImpl struct {
	userStorage UserStorage
	logger      *zap.SugaredLogger
}

func NewUserService(userStorage UserStorage, logs *logs.Logs) UserService {
	loggerName := "userService"
	logger := logs.WithName(loggerName)

	return &userServiceImpl{userStorage: userStorage, logger: logger}
}

func (service *userServiceImpl) Create(username, password string) (*User, error) {
	service.logger.Infow("Creating a user")

	hashedPassword, err := xpassword.Hash(password)
	if err != nil {
		service.logger.Errorf("Failed to hash password: %v", err)
		return nil, err
	}

	user, err := service.userStorage.Create(username, hashedPassword)
	if err != nil {
		service.logger.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	service.logger.Infow("User created", "id", user.Id, "username", user.Username)

	return user, nil
}

func (service *userServiceImpl) Get(userId domain.UserId) (*User, error) {
	user, err := service.userStorage.Get(userId)
	if err != nil {
		service.logger.Errorf("Failed to get user with id '%d': %v", userId, err)
		return nil, err
	}
	return user, nil
}

func (service *userServiceImpl) GetByUsername(username string) (*User, error) {
	user, err := service.userStorage.GetByUsername(username)
	if err != nil {
		service.logger.Errorf("Failed to get user by username '%s': %v", username, err)
		return nil, err
	}

	return user, nil
}

func (service *userServiceImpl) FindByUsername(username string) error {
	err := service.userStorage.FindByUsername(username)
	if err != nil {
		service.logger.Errorf("Failed to find user by username '%s': %v", username, err)

		return err
	}

	return nil
}
