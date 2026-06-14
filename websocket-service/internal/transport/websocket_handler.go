package transport

import (
	"fmt"
	"net/http"

	"github.com/excius/edns/websocket-service/internal/client"
	"github.com/excius/edns/websocket-service/internal/hub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Cors Origin check
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub *hub.Hub
}

func NewWebSocketHandler(h *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: h,
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

	if err := w.hub.SendToClient(client.ID, []byte("Welcome to EDNS!")); err != nil {
		fmt.Println("Failed to send to client")
		return
	}

	defer func() {
		w.hub.Unregister(client.ID)
		conn.Close()
	}()

	for {
		msgType, buff, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("message_type: ", msgType, ", buff: ", string(buff))
	}
}
