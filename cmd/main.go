package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found or error loading it. Relying on system environment variables.")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}
	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: connStr,
		},
	}

	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	api := application{
		config: cfg,
		db:     dbPool,
	}

	slog.Info("Starting server", "addr", cfg.addr)
	err = api.run(api.mount())
	if err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
