package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) UploadOrder(ctx context.Context, userID int, order *modelsOrder.OrderAccrual) error {
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

	queryInsert := `INSERT INTO orders (number, status, accrual, user_id) VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(ctx, queryInsert, order.Order, order.Status, order.Accrual, userID); err != nil {
		logger.Logger.Errorw("error inserting order", "error", err)

		return err
	}

	queryUpdateBalance := `UPDATE users SET current = (SELECT current FROM users WHERE id = $1) + $2 WHERE id = $1`
	result, err := tx.Exec(ctx, queryUpdateBalance, userID, order.Accrual)
	if err != nil {
		logger.Logger.Errorw("error updating balance", "error", err)

		return err
	}

	if result.RowsAffected() == 0 {
		logger.Logger.Errorw("no rows affected")

		return errors.New("error updating balance")
	}

	return nil
}

func (r *Repository) GetUserIdByOrder(ctx context.Context, order string) (userID int, err error) {
	query := `SELECT user_id FROM orders WHERE number = $1`

	if err := r.db.QueryRow(ctx, query, order).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		logger.Logger.Errorw("error in getting orders", "userID", userID, "err", err)

		return 0, err
	}

	return userID, err
}

func (r *Repository) GetOrders(ctx context.Context, userID int) ([]*modelsOrder.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1`
	var orders []*modelsOrder.Order

	// use null int because accrual can be NULL
	var accrual sql.NullFloat64

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		logger.Logger.Errorw("error in getting orders", "userID", userID, "err", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order modelsOrder.Order
		if err := rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt); err != nil {
			logger.Logger.Errorw("error in getting orders", "userID", userID, "err", err)

			return nil, err
		}

		// check accrual for not NULL value
		if accrual.Valid {
			order.Accrual = float32(accrual.Float64)
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, order, status string) error {
	result, err := r.db.Exec(ctx, `UPDATE orders SET status=$1 WHERE number=$2`, status, order)
	if err != nil {
		logger.Logger.Errorw("error updating order status", "error", err)

		return err
	}

	if ra := int(result.RowsAffected()); ra == 0 {
		logger.Logger.Errorw("no rows affected")

		return fmt.Errorf("no rows affected")
	}

	return nil
}
