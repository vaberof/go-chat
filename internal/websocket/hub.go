package websocket

import (
	"encoding/json"
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type Hub struct {
	clients    map[domain.UserId]*Client
	broadcast  chan MessageChan
	register   chan *Client
	unregister chan *Client

	messageService message.MessageService
	roomService    room.RoomService
	logger         *zap.SugaredLogger
}

func NewHub(roomService room.RoomService, logs *logs.Logs) *Hub {
	loggerName := "hub"
	logger := logs.WithName(loggerName)

	return &Hub{
		clients:     make(map[domain.UserId]*Client),
		broadcast:   make(chan MessageChan),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		roomService: roomService,
		logger:      logger,
	}
}

type Message struct {
	RoomId  domain.RoomId
	Message string
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case messageChan := <-h.broadcast:
			h.handleMessage(&messageChan)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client.userId] = client
	h.logger.Infof("Client with id: %d connected", client.userId)
}

func (h *Hub) unregisterClient(client *Client) {
	_, ok := h.clients[client.userId]
	if ok {
		delete(h.clients, client.userId)
		close(client.send)
	}
	h.logger.Infof("Client with id: %d disconnected", client.userId)
}

func (h *Hub) handleMessage(messageChan *MessageChan) {
	var message Message

	err := json.Unmarshal(messageChan.Message, &message)
	if err != nil {
		// send error
		return
	}

}
