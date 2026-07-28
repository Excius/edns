package app

import (
	"github.com/excius/edns/api-service/internal/handlers"
	"github.com/excius/edns/api-service/internal/metrics"
	"github.com/excius/edns/api-service/internal/middleware"
	"github.com/excius/edns/api-service/internal/service"
	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/observability/health"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/internal/stream"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router *gin.Engine
}

func NewApp(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client, metrics *metrics.Metrics) *App {

	// Repositories
	userRepo := repository.NewUserRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	// Queue
	stream := stream.NewRedisStream(redisClient, cfg.Redis.Stream)

	// Services
	userService := service.NewUserService(userRepo, notificationRepo, metrics)
	notificationService := service.NewNotificationService(userRepo, notificationRepo, deliveryRepo, stream, metrics)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Health
	readyHandler := health.NewHandler(
		health.NewDatabaseChecker(db),
		health.NewRedisChecker(redisClient),
	)

	// Router
	r := gin.New()

	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Metrics(metrics),
	)

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api.GET("/ready", func(c *gin.Context) {
		readyHandler.Ready(c)
	})

	api.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
