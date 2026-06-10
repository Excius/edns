package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/queue"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/worker-service/internal/processor"
	"go.uber.org/zap"
)

func main() {
	var wg sync.WaitGroup
	cfg := config.LoadConfig()

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db, err := config.NewPostgresPool(cfg)
	if err != nil {
		logger.Log.Error("DB connection failed:", zap.Error(err))
	}

	logger.Log.Info("Successfully connected to the database")
	defer db.Close()

	redisClient := config.NewRedisClient(cfg)
	defer redisClient.Close()

	// TODO: Need to make the worker name dynamic
	consumer := queue.NewRedisConsumer(redisClient, cfg.Redis.Stream, cfg.Redis.Group, "worker-1")
	revovery := queue.NewRedisRecovery(redisClient, cfg.Redis.Stream, cfg.Redis.Group, "worker-2", cfg.Redis.RecoveryInterval, cfg.Redis.RecoveryIdleTime)

	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	notificationProcessor := processor.NewNotificationProcessor(notificationRepo, deliveryRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = consumer.EnsureGroup(ctx); err != nil {
		logger.Log.Fatal(
			"Failed to create consumer group",
			zap.Error(err),
		)
		return
	}

	logger.Log.Info(
		"Consumer group ready",
		zap.String("group", cfg.Redis.Group),
	)

	wg.Go(func() {
		logger.Log.Info("Worker service is running...")
		if err := consumer.Start(ctx, notificationProcessor); err != nil && !errors.Is(err, context.Canceled) {
			logger.Log.Fatal("Worker failed", zap.Error(err))
		}
	})

	wg.Go(func() {
		if err := revovery.Start(ctx, notificationProcessor); err != nil && !errors.Is(err, context.Canceled) {
			logger.Log.Fatal("Recovery worker failed", zap.Error(err))
		}
	})

	// Listen for shutdown singnals
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Log.Info("Shutting down worker service...")

	cancel()

	wg.Wait()

	logger.Log.Info("Worker service stopped")
}
