package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"notification-system/api/internal/config"
	"notification-system/api/pkg/logger"
)

func main() {

	cfg := config.LoadConfig()

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db := config.NewPostgresPool(cfg)
	defer db.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	logger.Log.Info("Server started", zap.String("port", cfg.Server.Port))

	r.Run(":" + cfg.Server.Port)

}
