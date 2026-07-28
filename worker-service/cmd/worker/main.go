package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/internal/stream"
	"github.com/excius/edns/worker-service/internal/app"
	"github.com/excius/edns/worker-service/internal/metrics"
	"github.com/excius/edns/worker-service/internal/processor"
	"github.com/excius/edns/worker-service/internal/queue"
	"github.com/excius/edns/worker-service/internal/sender"
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
		return
	}

	logger.Log.Info("Successfully connected to the database")
	defer db.Close()

	mailClient, err := config.NewMailClient(cfg)
	if err != nil {
		logger.Log.Error("Mail client connection failed:", zap.Error(err))
		return
	}
	mailClient.Close()

	redisClient := config.NewRedisClient(cfg)
	defer redisClient.Close()

	metrics := metrics.NewMetrics()

	application := app.NewApp(db, redisClient, metrics)

	server := &http.Server{
		Addr:    ":" + cfg.WorkerServer.Port,
		Handler: application.Router,
	}

	hostname, _ := os.Hostname()
	consumerName := hostname

	consumer := queue.NewRedisConsumer(
		redisClient,
		cfg.Redis.Stream,
		cfg.Redis.Group,
		consumerName,
		&metrics.Consumer,
	)
	revovery := queue.NewRedisRecovery(
		redisClient,
		cfg.Redis.Stream,
		cfg.Redis.Group,
		consumerName,
		cfg.Redis.RecoveryMessageCount,
		cfg.Redis.RecoveryInterval,
		cfg.Redis.RecoveryIdleTime,
		&metrics.Recovery,
	)

	dlqPublisher := stream.NewRedisDLQStream(
		redisClient,
		cfg.Redis.DlqStream,
	)

	userRepo := repository.NewUserRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	deliveryRepo := repository.NewNotificationDeliveryRepository(db)

	senders := map[string]sender.Sender{
		string(models.ChannelEmail):     sender.NewEmailSender(mailClient, cfg.SMTP.From),
		string(models.ChannelWebsocket): sender.NewWebSocketSender(redisClient, cfg.Redis.Channel),
	}

	notificationProcessor := processor.NewNotificationProcessor(
		userRepo,
		notificationRepo,
		deliveryRepo,
		dlqPublisher,
		senders,
		&metrics.Processor,
	)

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
		logger.Log.Info("Starting worker server", zap.String("port", cfg.WorkerServer.Port))

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(
				"Worker server failed",
				zap.Error(err),
			)
		}
	})

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
