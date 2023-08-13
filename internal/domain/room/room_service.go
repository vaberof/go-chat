package room

import (
	"fmt"
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type RoomService interface {
	Create(creatorId domain.UserId, name, roomType string, members []domain.UserId) (*Room, error)
	Get(roomId domain.RoomId) (*Room, error)
	GetMembers(roomId domain.RoomId) ([]*Member, error)
	GetMessages(roomId domain.RoomId) ([]*message.Message, error)
	List(userId domain.UserId) ([]*Room, error)
}

type roomServiceImpl struct {
	roomStorage RoomStorage

	logger *zap.SugaredLogger
}

func NewRoomService(roomStorage RoomStorage, logs *logs.Logs) RoomService {
	loggerName := "room-service"
	logger := logs.WithName(loggerName)

	return &roomServiceImpl{roomStorage: roomStorage, logger: logger}
}

func (service *roomServiceImpl) Create(creatorId domain.UserId, name, roomType string, members []domain.UserId) (*Room, error) {
	if !service.isValidType(roomType) {
		service.logger.Errorf("Failed to create a room: invalid room type: '%s'", roomType)

		return nil, fmt.Errorf("'%s' is invalid room type", roomType)
	}

	room, err := service.roomStorage.Create(creatorId, name, roomType, members)
	if err != nil {
		service.logger.Errorf("Failed to create a room: %v", err)

		return nil, err
	}

	service.logger.Infow("Room created")

	return room, nil
}

func (service *roomServiceImpl) Get(roomId domain.RoomId) (*Room, error) {
	room, err := service.roomStorage.Get(roomId)
	if err != nil {
		service.logger.Errorf("Failed to get a room: %v", err)

		return nil, err
	}

	return room, nil
}

func (service *roomServiceImpl) GetMembers(roomId domain.RoomId) ([]*Member, error) {
	if err := service.roomStorage.Find(roomId); err != nil {
		return nil, err
	}

	return service.roomStorage.GetMembers(roomId)
}

func (service *roomServiceImpl) GetMessages(roomId domain.RoomId) ([]*message.Message, error) {
	if err := service.roomStorage.Find(roomId); err != nil {
		return nil, err
	}

	return service.roomStorage.GetMessages(roomId)
}

func (service *roomServiceImpl) List(userId domain.UserId) ([]*Room, error) {
	return service.roomStorage.List(userId)
}

func (service *roomServiceImpl) isValidType(roomType string) bool {
	return roomType == PrivateRoomType || roomType == GroupRoomType
}
