package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/observability/health"
	"github.com/excius/edns/websocket-service/internal/handler"
	"github.com/excius/edns/websocket-service/internal/hub"
	"github.com/excius/edns/websocket-service/internal/metrics"
	"github.com/excius/edns/websocket-service/internal/middleware"
	"github.com/excius/edns/websocket-service/internal/subscriber"
	"github.com/excius/edns/websocket-service/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Configuration
	cfg := config.LoadConfig()

	// Logger
	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	// Shared application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Redis
	redisClient := config.NewRedisClient(cfg)
	defer redisClient.Close()

	// Metrics
	metrics := metrics.NewMetrics()

	// Health
	readyHandler := health.NewHandler(
		health.NewRedisChecker(redisClient),
	)

	// Dependencies
	hub := hub.NewHub(metrics)

	notificationHandler := handler.NewNotificationHandler(hub)

	redisSubscriber := subscriber.NewRedisSubscriber(
		redisClient,
		cfg.Redis.Channel,
		metrics,
	)

	// Background workers
	var wg sync.WaitGroup

	wg.Go(func() {
		logger.Log.Info(
			"Starting Redis subscriber",
			zap.String("channel", cfg.Redis.Channel),
		)

		if err := redisSubscriber.Start(
			ctx,
			notificationHandler,
		); err != nil && !errors.Is(err, context.Canceled) {

			logger.Log.Error(
				"Redis subscriber failed",
				zap.Error(err),
			)

			cancel()
		}
	})

	// HTTP / WebSocket server
	wsHandler := transport.NewWebSocketHandler(hub, metrics)

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

	api.GET("/ws", wsHandler.Handler)

	server := &http.Server{
		Addr:    ":" + cfg.WSServer.Port,
		Handler: router,
	}

	go func() {
		logger.Log.Info(
			"Starting WebSocket server",
			zap.String("port", cfg.WSServer.Port),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			logger.Log.Fatal(
				"WebSocket server failed",
				zap.Error(err),
			)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logger.Log.Info("Shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error(
			"Server forced to shutdown",
			zap.Error(err),
		)
	}

	wg.Wait()

	logger.Log.Info("Server exited gracefully")
}
