package postgres

import (
	"errors"
	"fmt"
	"github.com/vaberof/go-chat/internal/domain/user"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	Rooms     []Room    `gorm:"many2many:members"`
	Messages  []Message `gorm:"foreignKey:SenderId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type userStorageImpl struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewUserStorage(db *gorm.DB, logs *logs.Logs) user.UserStorage {
	loggerName := "user-storage"
	logger := logs.WithName(loggerName)

	return &userStorageImpl{db: db, logger: logger}
}

func (storage *userStorageImpl) Create(username, password string) (*user.User, error) {
	var postgresUser User

	postgresUser.Username = username
	postgresUser.Password = password
	postgresUser.Rooms = make([]Room, 0)

	err := storage.db.Table("users").Create(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to create a user: %v", err)
		return nil, err
	}

	storage.logger.Infow("User created")

	return buildDomainUser(&postgresUser), nil
}

func (storage *userStorageImpl) Get(userId domain.UserId) (*user.User, error) {
	var postgresUser User

	err := storage.db.Preload("Rooms").Table("users").Where("id = ?", userId).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to get a user with id '%d': %v", userId, err)
		return nil, err
	}

	return buildDomainUser(&postgresUser), nil
}

func (storage *userStorageImpl) GetUsers(userIds []domain.UserId) ([]*user.User, error) {
	var postgresUsers []*User

	err := storage.db.Table("users").Where("id IN(?)", userIds).Find(&postgresUsers).Error
	if err != nil {
		storage.logger.Errorf("Failed to get users: %v", err)
		return nil, err
	}

	return buildDomainUsers(postgresUsers), nil
}

func (storage *userStorageImpl) GetByUsername(username string) (*user.User, error) {
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

func (storage *userStorageImpl) FindByUsername(username string) error {
	var postgresUser User

	err := storage.db.Table("users").Where("username = ?", username).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to find a user with username '%s': %v", username, err)
		return err
	}

	return nil
}

func buildDomainUser(postgresUser *User) *user.User {
	return &user.User{
		Id:       domain.UserId(postgresUser.Id),
		Username: postgresUser.Username,
		Password: postgresUser.Password,
		Rooms:    getRoomIds(postgresUser.Rooms),
		Messages: getMessagesIds(postgresUser.Messages),
	}
}

func buildDomainUsers(postgresUsers []*User) []*user.User {
	domainUsers := make([]*user.User, len(postgresUsers))

	for i := 0; i < len(domainUsers); i++ {
		domainUsers[i] = buildDomainUser(postgresUsers[i])
	}

	return domainUsers
}

func getRoomIds(rooms []Room) []domain.RoomId {
	roomIds := make([]domain.RoomId, len(rooms))

	for i := 0; i < len(roomIds); i++ {
		roomIds[i] = domain.RoomId(rooms[i].Id)
	}

	return roomIds
}
