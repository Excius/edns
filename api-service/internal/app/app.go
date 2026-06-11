package app

import (
	"github.com/excius/edns/api-service/internal/handlers"
	"github.com/excius/edns/api-service/internal/service"
	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/queue"
	"github.com/excius/edns/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router *gin.Engine
}

func NewApp(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client) *App {

	// Repositories
	userRepo := repository.NewUserRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	// Queue
	stream := queue.NewRedisStream(redisClient, cfg.Redis.Stream)

	// Services
	userService := service.NewUserService(userRepo, notificationRepo)
	notificationService := service.NewNotificationService(userRepo, notificationRepo, deliveryRepo, stream)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Router
	r := gin.Default()

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api.GET("/users/:id", userHandler.GetUserByID)
	api.GET("/users/email/:email", userHandler.GetUserByEmail)
	api.GET("/users/:id/notifications", userHandler.GetUserNotifications)
	api.POST("/users", userHandler.CreateUser)

	api.GET("/notifications/:id", notificationHandler.GetNotificationByID)
	api.GET("/notifications/:id/deliveries", notificationHandler.GetDeliveriesByNotificationID)
	api.POST("/notifications", notificationHandler.CreateNotification)

	return &App{
		Router: r,
	}
}

// TODO: Need to add rate limiting middleware to prevent abuse of the API.
