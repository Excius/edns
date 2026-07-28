package hub

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/websocket-service/internal/client"
	"github.com/excius/edns/websocket-service/internal/metrics"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Hub struct {
	// userID -> ConnID -> Conn
	users   map[string]map[string]struct{}
	clients map[string]*client.Client

	mu sync.RWMutex

	metrics *metrics.Metrics
}

func NewHub(metrics *metrics.Metrics) *Hub {
	return &Hub{
		users:   make(map[string]map[string]struct{}),
		clients: make(map[string]*client.Client),
		metrics: metrics,
	}
}

func (h *Hub) Register(c *client.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c.ID] = c

	if _, exists := h.users[c.UserID]; !exists {
		h.users[c.UserID] = make(map[string]struct{})
	}
	h.users[c.UserID][c.ID] = struct{}{}

	logger.Log.Info(
		"Client connected",
		zap.Int("total_clients", len(h.clients)),
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

	logger.Log.Info(
		"Client diconnected",
		zap.Int("total_clients", len(h.clients)),
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
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return errors.New("no conn associated")
	}

	err := client.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

func (h *Hub) SendToUser(userID string, payload []byte) error {

	h.mu.RLock()

	clientIDs := make([]string, 0, len(h.users[userID]))
	for clientID := range h.users[userID] {
		clientIDs = append(clientIDs, clientID)
	}

	h.mu.RUnlock()

	var sendErr error

	for _, clientID := range clientIDs {

		client := h.clients[clientID]

		start := time.Now()

		err := client.Conn.WriteMessage(websocket.TextMessage, payload)

		h.metrics.Delivery.DeliveryDuration.Observe(
			time.Since(start).Seconds(),
		)

		if err != nil {
			h.metrics.Delivery.DeliveryErrors.Inc()
			sendErr = fmt.Errorf("write failed: %w", err)
			continue
		}

		h.metrics.Delivery.MessagesDelivered.Inc()
	}

	return sendErr
}
