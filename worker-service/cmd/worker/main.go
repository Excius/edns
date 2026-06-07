package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/queue"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/worker-service/internal/processor"
	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db := config.NewPostgresPool(cfg, logger.Log)
	defer db.Close()

	redisClient := config.NewRedisClient(cfg)

	consumer := queue.NewRedisConsumer(redisClient, cfg.Redis.Stream)

	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	notificationProcessor := processor.NewNotificationProcessor(notificationRepo, deliveryRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Log.Info("Worker service is running...")

		if err := consumer.Start(ctx, notificationProcessor); err != nil && !errors.Is(err, context.Canceled) {
			logger.Log.Fatal("Worker failed", zap.Error(err))
		}
	}()

	// Listen for shutdown singnals
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Log.Info("Shutting down worker service...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	<-shutdownCtx.Done()

	logger.Log.Info("Worker service stopped")
}
