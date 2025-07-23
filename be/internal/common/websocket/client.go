package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// In production, you should restrict this
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// User ID for this client
	UserId uint64

	// Channels the user is subscribed to
	Channels map[uint64]bool
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Handle incoming messages (ping, typing indicators, etc.)
		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err == nil {
			c.handleIncomingMessage(&wsMsg)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
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

// handleIncomingMessage processes messages received from the client
func (c *Client) handleIncomingMessage(msg *WebSocketMessage) {
	switch msg.Type {
	case "ping":
		// Respond with pong
		pongMsg := &WebSocketMessage{
			Type:      "pong",
			Timestamp: uint64(time.Now().UnixMilli()),
		}
		if data, err := json.Marshal(pongMsg); err == nil {
			select {
			case c.send <- data:
			default:
				close(c.send)
			}
		}
	case "typing":
		// Broadcast typing indicator to channel members
		if msg.ChannelId > 0 {
			typingMsg := &WebSocketMessage{
				Type:      "typing",
				ChannelId: msg.ChannelId,
				UserId:    c.UserId,
				Data:      msg.Data,
				Timestamp: uint64(time.Now().UnixMilli()),
			}
			c.hub.BroadcastToChannel(msg.ChannelId, typingMsg, c.UserId)
		}
	case "subscribe":
		// Subscribe to channel updates
		if channelId, ok := msg.Data.(float64); ok {
			if c.Channels == nil {
				c.Channels = make(map[uint64]bool)
			}
			c.Channels[uint64(channelId)] = true
			log.Printf("User %d subscribed to channel %d", c.UserId, uint64(channelId))
		}
	case "unsubscribe":
		// Unsubscribe from channel updates
		if channelId, ok := msg.Data.(float64); ok {
			if c.Channels != nil {
				delete(c.Channels, uint64(channelId))
			}
			log.Printf("User %d unsubscribed from channel %d", c.UserId, uint64(channelId))
		}
	}
}

// ServeWS handles websocket requests from the peer.
func ServeWS(hub *Hub, ctx *gin.Context) {
	// Extract user ID from authentication (already set by auth middleware)
	userIdInterface, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId, ok := userIdInterface.(uint64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		UserId:   userId,
		Channels: make(map[uint64]bool),
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
