package PostgreSQL

import (
	"context"
)

func (r *Repository) Register(ctx context.Context, login, password string) error {
	qury := `INSERT INTO users (login, password) VALUES ($1, $2)`

	result, err := r.db.Exec(ctx, login, password)
	if err != nil {
		logger.Logger

		return err
	}
}
