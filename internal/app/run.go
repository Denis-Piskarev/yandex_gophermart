package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/config"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	_ "github.com/DenisquaP/yandex_gophermart/migrations"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func init() {
	logger.NewLogger()
}

// Run - starts app
func Run() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Logger.Fatalw("Failed to get config", "error", err)
	}

	conn, err := pgx.Connect(ctx, cfg.DatabaseUri)
	if err != nil {
		logger.Logger.Fatalw("Failed to connect to database", "error", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			logger.Logger.Fatalw("Failed to close connection", "error", err)
		}
	}()

	if err := migrate(cfg.DatabaseUri); err != nil {
		logger.Logger.Fatalw("Failed to migrate database", "error", err)

		return
	}

	runServer(cfg)
}

// runServer - starts server
func runServer(cfg *config.Config) {
	logger.Logger.Infow(fmt.Sprintf("Starting server on %s...", cfg))
}

// Migrates to database
func migrate(addr string) error {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		logger.Logger.Errorw("Failed to open DB", "error", err)

		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Logger.Errorw("Failed to close connection", "error", err)
		}
	}()

	if err := goose.Up(db, "./migrations"); err != nil {
		logger.Logger.Errorw("Failed to run migrations", "error", err)

		return err
	}

	return nil
}
