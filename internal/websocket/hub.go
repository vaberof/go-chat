package websocket

import (
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/internal/domain/user"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
)

type Hub struct {
	clients map[domain.UserId]*Client
	rooms   map[domain.RoomId]*Room

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	messageService message.MessageService
	roomService    room.RoomService
	userService    user.UserService
	logger         *zap.SugaredLogger
}

func NewHub(messageService message.MessageService, roomService room.RoomService, userService user.UserService, logs *logs.Logs) *Hub {
	loggerName := "hub"
	logger := logs.WithName(loggerName)

	return &Hub{
		clients:        make(map[domain.UserId]*Client),
		rooms:          make(map[domain.RoomId]*Room),
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		messageService: messageService,
		roomService:    roomService,
		userService:    userService,
		logger:         logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToClients(message)
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

func (h *Hub) broadcastToClients(message []byte) {
	for _, client := range h.clients {
		client.send <- message
	}
}

func (h *Hub) broadcastToRoom(action string, message *MessagePayload) {
	roomId := message.RoomId

	room, ok := h.rooms[roomId]
	if !ok {
		room = h.createRoom(roomId)
	}

	room.broadcast <- &Message{
		Action:  action,
		Payload: message.Encode(),
	}
}

func (h *Hub) createRoom(roomId domain.RoomId) *Room {
	room := NewRoom(roomId, h.userService)
	h.rooms[roomId] = room
	go room.Run()

	return room
}
