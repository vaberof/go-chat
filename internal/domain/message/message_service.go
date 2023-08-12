package message

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type MessageService interface {
	Create(senderId domain.UserId, roomId domain.RoomId, message string) (*Message, error)
}

type messageServiceImpl struct {
	messageStorage MessageStorage

	logger *zap.SugaredLogger
}

func NewMessageService(messageStorage MessageStorage, logs *logs.Logs) MessageService {
	loggerName := "message-service"
	logger := logs.WithName(loggerName)

	return &messageServiceImpl{
		messageStorage: messageStorage, logger: logger,
	}
}

func (service *messageServiceImpl) Create(senderId domain.UserId, roomId domain.RoomId, message string) (*Message, error) {
	return nil, nil
}
