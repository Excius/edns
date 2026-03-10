package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewPostgres() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://admin:admin@localhost:5432/notifications")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	log.Println("Connected to PostgreSQL")

	return conn
}
