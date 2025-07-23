package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// User to clients mapping for targeted messaging
	userClients map[uint64][]*Client
	mutex       sync.RWMutex
}

// WebSocketMessage represents the structure of messages sent over WebSocket
type WebSocketMessage struct {
	Type      string      `json:"type"` // "message", "presence", "membership", "typing"
	Data      interface{} `json:"data"` // The actual message data
	ChannelId uint64      `json:"channel_id,omitempty"`
	UserId    uint64      `json:"user_id,omitempty"`
	Timestamp uint64      `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[uint64][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			if client.UserId > 0 {
				h.userClients[client.UserId] = append(h.userClients[client.UserId], client)
			}
			h.mutex.Unlock()
			log.Printf("Client registered: userId=%d, total clients=%d", client.UserId, len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from user clients mapping
				if client.UserId > 0 {
					userClients := h.userClients[client.UserId]
					for i, c := range userClients {
						if c == client {
							h.userClients[client.UserId] = append(userClients[:i], userClients[i+1:]...)
							break
						}
					}
					if len(h.userClients[client.UserId]) == 0 {
						delete(h.userClients, client.UserId)
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: userId=%d, total clients=%d", client.UserId, len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// BroadcastToChannel sends a message to all users in a specific channel
func (h *Hub) BroadcastToChannel(channelId uint64, message *WebSocketMessage, excludeUserId uint64) {
	// This is a simplified version. In a production system, you'd want to:
	// 1. Query the database for channel members
	// 2. Send only to those members
	// For now, we'll broadcast to all connected clients except the sender

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	h.mutex.RLock()
	for client := range h.clients {
		// Skip the sender
		if client.UserId == excludeUserId {
			continue
		}

		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
	h.mutex.RUnlock()
}

// SendToUser sends a message to all connections of a specific user
func (h *Hub) SendToUser(userId uint64, message *WebSocketMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	h.mutex.RLock()
	clients := h.userClients[userId]
	for _, client := range clients {
		select {
		case client.send <- messageBytes:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
	h.mutex.RUnlock()
}

// SendToUsers sends a message to multiple users
func (h *Hub) SendToUsers(userIds []uint64, message *WebSocketMessage) {
	for _, userId := range userIds {
		h.SendToUser(userId, message)
	}
}

// GetConnectedUsers returns a list of currently connected user IDs
func (h *Hub) GetConnectedUsers() []uint64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	userIds := make([]uint64, 0, len(h.userClients))
	for userId := range h.userClients {
		userIds = append(userIds, userId)
	}
	return userIds
}

// IsUserOnline checks if a user has any active connections
func (h *Hub) IsUserOnline(userId uint64) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	clients, exists := h.userClients[userId]
	return exists && len(clients) > 0
}
