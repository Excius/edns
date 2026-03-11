package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"notification-system/api/internal/config"
	"notification-system/api/internal/handlers"
	"notification-system/api/internal/repository"
	"notification-system/api/internal/service"
	"notification-system/api/pkg/logger"
)

func main() {

	cfg := config.LoadConfig()

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db := config.NewPostgresPool(cfg)
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	// Services
	userService := service.NewUserService(userRepo)
	notificationService := service.NewNotificationService(userRepo, notificationRepo, deliveryRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Router
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.POST("/users", userHandler.CreateUser)
	api.POST("/notifications", notificationHandler.CreateNotification)

	logger.Log.Info("Server started", zap.String("port", cfg.Server.Port))

	r.Run(":" + cfg.Server.Port)

}
