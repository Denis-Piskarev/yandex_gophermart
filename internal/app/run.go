package app

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/config"
	_ "github.com/DenisquaP/yandex_gophermart/migrations"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

// Run - starts app
func Run() {
	ctx := context.Background()

	newDev, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := newDev.Sync(); err != nil {
			log.Println(fmt.Errorf("failed to sync %w", err))
		}
	}()

	logger := newDev.Sugar()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("Failed to get config", "error", err)
	}

	conn, err := pgx.Connect(ctx, cfg.DatabaseUri)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "error", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			logger.Fatalw("Failed to close connection", "error", err)
		}
	}()

	if err := migrate(logger, cfg.DatabaseUri); err != nil {
		logger.Fatalw("Failed to migrate database", "error", err)

		return
	}

	runServer(cfg, logger)
}

// runServer - starts server
func runServer(cfg *config.Config, logger *zap.SugaredLogger) {
	logger.Infow(fmt.Sprintf("Starting server on %s...", cfg))
}

//go:embed *.go
var embedMigrations embed.FS

// Migrates to database
func migrate(logger *zap.SugaredLogger, addr string) error {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		logger.Errorw("Failed to open DB", "error", err)

		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorw("Failed to close connection", "error", err)
		}
	}()

	if err := goose.Up(db, "./migrations"); err != nil {
		logger.Errorw("Failed to run migrations", "error", err)

		return err
	}

	return nil
}
