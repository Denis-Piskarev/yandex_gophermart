package postgresql

import (
	"context"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func (r *Repository) Withdraw(ctx context.Context, userID int, sum float32, order int) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Logger.Errorw("Error starting transaction", "error", err)

		return err
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(ctx); txErr != nil {
				logger.Logger.Errorw("Error rolling back", "error", txErr)

				return
			}
		}
		if err = tx.Commit(ctx); err != nil {
			logger.Logger.Errorw("Error committing transaction", "error", err)

		}
	}()

	var current float32
	queryGetCurrent := `SELECT current FROM users WHERE id = $1`
	if err := tx.QueryRow(ctx, queryGetCurrent, userID).Scan(&current); err != nil {
		logger.Logger.Errorw("Error getting current balance of user", "userID", userID, "error", err)

		return err
	}

	// return err if not enough balance
	if sum > current {
		logger.Logger.Errorw("not enough balance for withdraw", "userID", userID, "sum", sum, "current", current)
		cErr := customerrors.NewCustomError("not enough balance", http.StatusPaymentRequired)

		return cErr
	}

	// adding withdraw to user
	queryAddWithdrawToUser := `UPDATE users SET withdrawn = (SELECT withdrawn FROM users WHERE id = $1) + $2 WHERE id = $1`
	if _, err := tx.Exec(ctx, queryAddWithdrawToUser, userID, sum); err != nil {
		logger.Logger.Errorw("Error adding withdraw", "userID", userID, "error", err)

		return err
	}

	// adding withdraw to table withdraws
	queryAddWithdrawToTable := `INSERT INTO withdrawals (user_id, number, sum) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(ctx, queryAddWithdrawToTable, userID, order, sum); err != nil {
		logger.Logger.Errorw("Error adding withdraw to table", "userID", userID, "error", err)

		return err
	}

	return nil
}
