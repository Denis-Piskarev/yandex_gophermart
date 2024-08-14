package PostgreSQL

import (
	"context"
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
)

func (r *Repository) Register(ctx context.Context, login, password string) error {
	query := `INSERT INTO users (login, password) VALUES ($1, $2)`

	result, err := r.db.Exec(ctx, query, login, password)
	if err != nil {
		logger.Logger.Errorw("error inserting user", "login", login, "error", err)

		return err
	}

	if result.RowsAffected() == 0 {
		logger.Logger.Errorw("error inserting user", "login", login, "error", "no rows inserted")

		return errors.New("no rows inserted")
	}

	return nil
}

func (r *Repository) CheckLogin(ctx context.Context, login string) error {
	var userId int

	if err := r.db.QueryRow(ctx, `SELECT userId FROM users WHERE login=$1`, login).Scan(&userId); err != nil {
		logger.Logger.Errorw("error querying user", "login", login, "error", err)

		return err
	}

	if userId != 0 {
		err := models.CustomError{
			Err:        "user already exists",
			StatusCode: http.StatusConflict,
		}

		logger.Logger.Errorw("user already exists", "login", login, "error", err.Error())

		return err
	}

	return nil
}
