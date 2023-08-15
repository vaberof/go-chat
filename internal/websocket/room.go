package websocket

import (
	"fmt"
	"github.com/vaberof/go-chat/internal/domain/user"
	"github.com/vaberof/go-chat/pkg/domain"
)

type Room struct {
	roomId domain.RoomId

	clients    map[domain.UserId]*Client
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client

	userService user.UserService
}

func NewRoom(roomId domain.RoomId, userService user.UserService) *Room {
	return &Room{
		roomId:      roomId,
		clients:     make(map[domain.UserId]*Client),
		broadcast:   make(chan *Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userService: userService,
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.registerClient(client)
		case client := <-r.unregister:
			r.unregisterClient(client)
		case message := <-r.broadcast:
			r.broadcastToClients(message)
		}
	}
}

func (r *Room) registerClient(client *Client) {
	r.clients[client.userId] = client
	r.notifyJoinRoom(client.userId)
}

func (r *Room) unregisterClient(client *Client) {
	_, ok := r.clients[client.userId]
	if ok {
		delete(r.clients, client.userId)
	}
}

func (r *Room) broadcastToClients(message *Message) {
	for _, client := range r.clients {
		client.send <- message.Encode()
	}
}

func (r *Room) notifyJoinRoom(userId domain.UserId) {
	user, err := r.userService.Get(userId)
	if err != nil {
		return
	}

	for _, client := range r.clients {
		client.send <- []byte(fmt.Sprintf("User %s joined to the room with id %d", user.Username, r.roomId))
	}
}
