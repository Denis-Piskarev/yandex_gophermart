package PostgreSQL

import (
	"context"
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customErrors"
	"github.com/jackc/pgx/v5"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
)

func (r *Repository) Register(ctx context.Context, login, password string) (int, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	var id int

	if err := r.db.QueryRow(ctx, query, login, password).Scan(&id); err != nil {
		logger.Logger.Error("error inserting user", "login", login, "error", err)

		return 0, err
	}

	return 0, nil
}

func (r *Repository) Login(ctx context.Context, login, password string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2;`
	var id int

	if err := r.db.QueryRow(ctx, query, login, password).Scan(&id); err != nil {
		logger.Logger.Errorw("error getting user", "login", login, "error", err)

		return 0, err
	}

	return id, nil
}

func (r *Repository) CheckLogin(ctx context.Context, login string) error {
	var userId int

	if err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE login=$1`, login).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		logger.Logger.Errorw("error querying user", " login ", login, " error ", err)

		return err
	}

	if userId != 0 {
		err := customErrors.CustomError{
			Err:        "user already exists",
			StatusCode: http.StatusConflict,
		}

		logger.Logger.Errorw("user already exists", "login", login, "error", err.Error())

		return err
	}

	return nil
}
