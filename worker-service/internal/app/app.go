package app

import (
	"github.com/excius/edns/internal/observability/health"
	"github.com/excius/edns/worker-service/internal/metrics"
	"github.com/excius/edns/worker-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router *gin.Engine
}

func NewApp(db *pgxpool.Pool, redisClient *redis.Client, metrics *metrics.Metrics) *App {

	readyHandler := health.NewHandler(
		health.NewDatabaseChecker(db),
		health.NewRedisChecker(redisClient),
	)

	// Router
	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
	)

	api := router.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api.GET("/ready", func(c *gin.Context) {
		readyHandler.Ready(c)
	})

	api.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return &App{router}
}
