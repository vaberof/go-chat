package websocket

import (
	"github.com/vaberof/go-chat/internal/domain/chat/room"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type Hub struct {
	clients   map[user.UserId]*Client
	chatRooms map[room.RoomId]*room.Room

	register   chan *Client
	unregister chan *Client

	logger *zap.SugaredLogger
}

func NewHub(logs *logs.Logs) *Hub {
	loggerName := "hub"
	logger := logs.WithName(loggerName)
	return &Hub{
		clients:    make(map[user.UserId]*Client),
		chatRooms:  make(map[room.RoomId]*room.Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}
