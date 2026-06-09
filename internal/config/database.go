package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FIX: Need to move the logger to err from the function
func NewPostgresPool(cfg *Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse database URL: %w", err)
	}

	config.MaxConns = cfg.DB.MAXConns
	config.MinConns = cfg.DB.MINConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return pool, nil
}
