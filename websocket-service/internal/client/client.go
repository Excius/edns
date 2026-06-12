package client

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
}
