package postgresql

import (
	"context"
	"errors"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Register(ctx context.Context, login, password string) (int, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	var id int

	if err := r.db.QueryRow(ctx, query, login, password).Scan(&id); err != nil {
		logger.Logger.Error("error inserting user", "login", login, "error", err)

		return 0, err
	}

	return id, nil
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

func (r *Repository) GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error) {
	var withdrawals []*modelsUser.Withdrawals

	query := `SELECT number, sum, processed_at FROM withdrawals WHERE user_id=$1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		logger.Logger.Errorw("error getting withdrawals", "error", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var w modelsUser.Withdrawals
		if err := rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			logger.Logger.Errorw("error scanning withdrawals", "error", err)

			return nil, err
		}

		withdrawals = append(withdrawals, &w)
	}

	// if user has no withdrawals return custom error with status code
	if len(withdrawals) == 0 {
		cErr := customerrors.NewCustomError("user has no withdrawals", http.StatusNoContent)

		return nil, cErr
	}

	return withdrawals, nil
}

func (r *Repository) CheckLogin(ctx context.Context, login string) error {
	var userID int

	if err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE login=$1`, login).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		logger.Logger.Errorw("error querying user", " login ", login, " error ", err)

		return err
	}

	if userID != 0 {
		err := customerrors.CustomError{
			Err:        "user already exists",
			StatusCode: http.StatusConflict,
		}

		logger.Logger.Errorw("user already exists", "login", login, "error", err.Error())

		return err
	}

	return nil
}
