package main

import (
	"context"
	"errors"
	"net/http"
	"notification-system/api/internal/app"
	"notification-system/api/internal/config"
	"notification-system/api/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	cfg := config.LoadConfig()

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db := config.NewPostgresPool(cfg)
	defer db.Close()

	application := app.NewApp(cfg, db)

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
