package postgres

import (
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	Id        int64 `gorm:"primaryKey"`
	SenderId  int64
	RoomId    int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type messageStorageImpl struct {
	db *gorm.DB

	logger *zap.SugaredLogger
}

func NewMessageStorage(db *gorm.DB, logs *logs.Logs) message.MessageStorage {
	loggerName := "message-storage"
	logger := logs.WithName(loggerName)

	return &messageStorageImpl{db: db, logger: logger}
}

func (storage *messageStorageImpl) Create(senderId domain.UserId, roomId domain.RoomId, message string) (*message.Message, error) {
	var postgresMessage Message

	postgresMessage.SenderId = int64(senderId)
	postgresMessage.RoomId = int64(roomId)
	postgresMessage.Text = message

	err := storage.db.Table("messages").Create(&postgresMessage).Error
	if err != nil {
		return nil, err
	}

	return buildDomainMessage(&postgresMessage), nil
}

func buildDomainMessage(postgresMessage *Message) *message.Message {
	return &message.Message{
		Id:        domain.MessageId(postgresMessage.Id),
		SenderId:  domain.UserId(postgresMessage.SenderId),
		RoomId:    domain.RoomId(postgresMessage.RoomId),
		Text:      postgresMessage.Text,
		CreatedAt: postgresMessage.CreatedAt,
		EditedAt:  postgresMessage.UpdatedAt,
	}
}

func buildDomainMessages(postgresMessages []*Message) []*message.Message {
	domainMessages := make([]*message.Message, len(postgresMessages))

	for i := 0; i < len(domainMessages); i++ {
		domainMessages[i] = buildDomainMessage(postgresMessages[i])
	}

	return domainMessages
}
