package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

// ConnectionDB устанавливает подключение к базе данных PostgreSQL
// и возвращает пул подключений для дальнейшего использования.
// Используются настройки из переменной окружения DATABASE_URL.
// Максимальное количество соединений - 10, максимальное время простоя - 5 минут.
func ConnectionDB() *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return pool
}
