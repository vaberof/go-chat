package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	userId user.UserId
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}

func NewClient(userId user.UserId, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{hub: hub, userId: userId, conn: conn, send: make(chan []byte, 256)}
}

func serveWebSocket(hub *Hub, w http.ResponseWriter, req *http.Request) {
}
