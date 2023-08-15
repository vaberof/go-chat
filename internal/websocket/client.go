package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/vaberof/go-chat/pkg/auth"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	userId domain.UserId
	rooms  map[domain.RoomId]*Room
	conn   *websocket.Conn
	hub    *Hub
	send   chan []byte

	logger *zap.SugaredLogger
}

func NewClient(userId domain.UserId, hub *Hub, conn *websocket.Conn, logs *logs.Logs) *Client {
	loggerName := "websocket-client"
	logger := logs.WithName(loggerName)

	return &Client{
		hub:    hub,
		userId: userId,
		rooms:  make(map[domain.RoomId]*Room),
		conn:   conn,
		send:   make(chan []byte, 256),
		logger: logger,
	}
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.handleMessage(jsonMessage)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request, logs *logs.Logs) {
	userId := auth.UserIdFromContext(r.Context())
	if userId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing jwt token"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(*userId, hub, conn, logs)
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) handleMessage(jsonMessage []byte) {
	var inboundMessage Message

	err := json.Unmarshal(jsonMessage, &inboundMessage)
	if err != nil {
		return
	}

	switch inboundMessage.Action {
	case SendMessageAction:
		var payload MessagePayload

		err := json.Unmarshal(inboundMessage.Payload, &payload)
		if err != nil {
			c.logger.Errorf("Failed to unmarshal to MessagePayload %v\n", err)
			return
		}

		c.handleSendMessageAction(inboundMessage.Action, &payload)

	case JoinRoomAction:
		var payload JoinRoomPayload

		err := json.Unmarshal(inboundMessage.Payload, &payload)
		if err != nil {
			c.logger.Errorf("Failed to unmarshal to JoinRoomPayload %v\n", err)
			return
		}

		c.handleJoinRoomAction(&payload)
	}
}

func (c *Client) handleSendMessageAction(action string, messagePayload *MessagePayload) {
	_, err := c.hub.messageService.Create(messagePayload.SenderId, messagePayload.RoomId, messagePayload.Text)
	if err != nil {
		c.logger.Errorf("Failed to create message: %v", err)
		return
	}

	c.hub.broadcastToRoom(action, messagePayload)
}

func (c *Client) handleJoinRoomAction(message *JoinRoomPayload) {
	roomId := message.RoomId

	_, err := c.hub.roomService.Get(roomId)
	if err != nil {
		c.logger.Errorf("Failed to get room with id %d: %v\n", roomId, err)
		return
	}

	room, ok := c.hub.rooms[roomId]
	if !ok {
		room = c.hub.createRoom(roomId)
	}

	c.rooms[roomId] = room

	room.register <- c
}
