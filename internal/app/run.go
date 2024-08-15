package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/config"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/repository/PostgreSQL"
	"github.com/DenisquaP/yandex_gophermart/internal/rest/endpoints"
	"github.com/DenisquaP/yandex_gophermart/internal/rest/router"
	"github.com/DenisquaP/yandex_gophermart/internal/service"
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
		return
	}

	repo := PostgreSQL.NewRepository(conn)
	serv := service.NewService(repo)
	endPoints := endpoints.NewEndpoints(serv)
	routes := router.NewRouterWithMiddleware(endPoints, serv)

	runServer(cfg, routes)
}

// runServer - starts server
func runServer(cfg *config.Config, handler http.Handler) {
	logger.Logger.Infow(fmt.Sprintf("Starting server on %s...", cfg))

	if err := http.ListenAndServe(cfg.RunAddress, handler); err != nil {
		logger.Logger.Fatalw("Failed to start server", "error", err)
	}
}

// Migrates to database
func migrate(addr string) error {
	db, err := sql.Open("pgx", addr)
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
