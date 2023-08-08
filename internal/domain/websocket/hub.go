package websocket

import (
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type Hub struct {
	clients    map[domain.UserId]*Client
	broadcast  chan MessageChan
	register   chan *Client
	unregister chan *Client

	logger *zap.SugaredLogger
}

func NewHub(logs *logs.Logs) *Hub {
	loggerName := "hub"
	logger := logs.WithName(loggerName)

	return &Hub{
		clients:    make(map[domain.UserId]*Client),
		broadcast:  make(chan MessageChan),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {

}
