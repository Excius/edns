package main

import (
	"github.com/excius/edns/websocket-service/internal/handlers"
	"github.com/excius/edns/websocket-service/internal/hub"
	"github.com/gin-gonic/gin"
)

func main() {

	hub := hub.NewHub()

	handler := handlers.NewWebSocketHandler(hub)

	router := gin.Default()

	router.GET("/ws", handler.Handler)

	router.Run(":8081")
}
