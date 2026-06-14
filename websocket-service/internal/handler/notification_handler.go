package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/excius/edns/internal/events"
	"github.com/excius/edns/websocket-service/internal/hub"
)

type WebSocketNotification struct {
	NotificationID string `json:"notification_id"`
	Title          string `json:"title"`
	Message        string `json:"message"`
}

type NotificationHandler struct {
	hub *hub.Hub
}

func NewNotificationHandler(hub *hub.Hub) *NotificationHandler {
	return &NotificationHandler{
		hub: hub,
	}
}

func (n *NotificationHandler) Handle(ctx context.Context, payload []byte) error {

	var event events.NotificationEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf(
			"unmarshel notification error: %w",
			err,
		)
	}

	if event.UserID == "" {
		return errors.New("missing user_id")
	}

	notification := WebSocketNotification{
		NotificationID: event.NotificationID,
		Title:          event.Title,
		Message:        event.Message,
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshel notification failed: %w", err)
	}

	return n.hub.SendToUser(
		event.UserID,
		payload,
	)
}
