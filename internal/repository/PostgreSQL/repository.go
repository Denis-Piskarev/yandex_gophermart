package PostgreSQL

import (
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository struct {
	db     *pgx.Conn
	logger *zap.SugaredLogger
}

func NewRepository(db *pgx.Conn, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
