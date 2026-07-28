package transport

import (
	"net/http"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/websocket-service/internal/client"
	"github.com/excius/edns/websocket-service/internal/hub"
	"github.com/excius/edns/websocket-service/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	// Cors Origin check
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub     *hub.Hub
	metrics *metrics.Metrics
}

func NewWebSocketHandler(h *hub.Hub, metrics *metrics.Metrics) *WebSocketHandler {
	return &WebSocketHandler{
		hub:     h,
		metrics: metrics,
	}
}

func (w *WebSocketHandler) Handler(c *gin.Context) {

	userID := c.Query("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &client.Client{
		ID:     uuid.NewString(),
		Conn:   conn,
		UserID: userID,
	}

	w.hub.Register(client)

	w.metrics.Connection.ConnectionsOpened.Inc()
	w.metrics.Connection.ActiveConnections.Inc()

	defer func() {
		w.metrics.Connection.ConnectionsClosed.Inc()
		w.metrics.Connection.ActiveConnections.Dec()

		w.hub.Unregister(client.ID)

		if err := conn.Close(); err != nil {
		}
	}()

	if err := w.hub.SendToClient(client.ID, []byte("Welcome to EDNS!")); err != nil {
		logger.Log.Info("Failed to send welcome message")
		return
	}

	for {
		msgType, buff, err := conn.ReadMessage()
		if err != nil {
			break
		}

		logger.Log.Info(
			"message info",
			zap.Int("message_type", msgType),
			zap.String("buff", string(buff)),
		)
	}
}
