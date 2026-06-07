package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPostgresPool(cfg *Config, log *zap.Logger) *pgxpool.Pool {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		log.Error("Unable to parse database URL:", zap.Error(err))
	}

	config.MaxConns = cfg.DB.MAXConns
	config.MinConns = cfg.DB.MINConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error("Unable to create connection pool:", zap.Error(err))
	}

	if err := pool.Ping(ctx); err != nil {
		log.Error("Unable to connect to database:", zap.Error(err))
	}

	log.Info("Successfully connected to the database")

	return pool
}
