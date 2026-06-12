package hub

import (
	"errors"
	"fmt"
	"sync"

	"github.com/excius/edns/websocket-service/internal/client"
	"github.com/gorilla/websocket"
)

type Hub struct {
	// userID -> ConnID -> Conn
	users   map[string]map[string]*client.Client
	clients map[string]*client.Client

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		users:   make(map[string]map[string]*client.Client),
		clients: make(map[string]*client.Client),
	}
}

func (h *Hub) Register(c *client.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c.ID] = c

	if _, exists := h.users[c.UserID]; !exists {
		h.users[c.UserID] = make(map[string]*client.Client)
	}
	h.users[c.UserID][c.ID] = c

	fmt.Println(
		"Client connected. Total clients:",
		len(h.clients),
	)
}

func (h *Hub) Unregister(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.clients[clientID]

	if exists {
		delete(h.clients, clientID)

		if userConnection, exists := h.users[conn.UserID]; exists {
			delete(h.users[conn.UserID], clientID)

			if len(userConnection) == 0 {
				delete(h.users, conn.UserID)
			}
		}
	}

	fmt.Println(
		"Client disconnected. Total clients:",
		len(h.clients),
	)
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	length := len(h.clients)
	return length
}

func (h *Hub) SendToClient(clientID string, message []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[clientID]
	if !ok {
		return errors.New("no conn associated")
	}

	err := client.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (h *Hub) SendToUser(userId string, message []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.users[userId]
	if !ok {
		return errors.New("no conn associated with this user")
	}

	var sendErr error

	for _, client := range conns {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			sendErr = fmt.Errorf("write failed: %w", err)
		}
	}
	return sendErr
}
