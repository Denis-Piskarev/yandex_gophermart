package postgresql

import (
	"context"
	"errors"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetBalance(ctx context.Context, userID int) (modelsBalance.Balance, error) {
	var balance modelsBalance.Balance

	query := `SELECT current, withdrawn FROM users WHERE id=$1`
	if err := r.db.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn); err != nil {
		// if user has no orders we return default values of balance
		if errors.Is(err, pgx.ErrNoRows) {
			return modelsBalance.Balance{}, nil
		}

		logger.Logger.Errorw("error getting user`s balance", "error", err)

		return modelsBalance.Balance{}, err
	}

	return balance, nil
}
