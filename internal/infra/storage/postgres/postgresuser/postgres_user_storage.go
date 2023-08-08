package postgresuser

import (
	"errors"
	"fmt"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserStorage struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewUserStorage(db *gorm.DB, logs *logs.Logs) *UserStorage {
	loggerName := "user-storage"
	logger := logs.WithName(loggerName)

	return &UserStorage{db: db, logger: logger}
}

func (storage *UserStorage) Create(username, password string) (*user.User, error) {
	var postgresUser User

	postgresUser.Username = username
	postgresUser.Password = password

	err := storage.db.Table("users").Create(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to create a user: %v", err)
		return nil, err
	}

	storage.logger.Infow("User created")

	return buildDomainUser(&postgresUser), nil
}

func (storage *UserStorage) Get(userId domain.UserId) (*user.User, error) {
	var postgresUser User

	err := storage.db.Table("users").Where("id = ?", userId).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to get a user with id '%d': %v", userId, err)
		return nil, err
	}

	return buildDomainUser(&postgresUser), nil
}

func (storage *UserStorage) GetByUsername(username string) (*user.User, error) {
	var postgresUser User

	err := storage.db.Table("users").Where("username = ?", username).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to get a user with username '%s': %v", username, err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User with username '%s' not found", username)
		}

		return nil, err
	}

	return buildDomainUser(&postgresUser), nil
}

func (storage *UserStorage) FindByUsername(username string) error {
	var postgresUser User

	err := storage.db.Table("users").Where("username = ?", username).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to find a user with username '%s': %v", username, err)
		return err
	}

	return nil
}
