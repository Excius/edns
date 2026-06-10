package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/excius/edns/api-service/internal/app"
	"github.com/excius/edns/internal/config"
	"github.com/excius/edns/internal/logger"
	"go.uber.org/zap"
)

func main() {

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

	application := app.NewApp(cfg, db, redisClient)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: application.Router,
	}

	go func() {
		logger.Log.Info("Starting server", zap.String("port", cfg.Server.Port))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Listen for shutdown signals and gracefully shut down the server
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")

}
