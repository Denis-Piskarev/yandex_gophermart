package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UploadOrder(ctx context.Context, userID int, order *modelsOrder.OrderAccrual) error {
	query := `INSERT INTO orders (number, status, accrual, user_id) VALUES ($1, $2, $3, $4)`

	orderInt, err := strconv.Atoi(order.Order)
	if err != nil {
		logger.Logger.Errorw("error converting order to integer", "error", err)

		return err
	}

	result, err := r.db.Exec(ctx, query, orderInt, order.Status, order.Accrual, userID)
	if err != nil {
		logger.Logger.Errorw("error inserting order", "error", err)

		return err
	}

	if result.RowsAffected() == 0 {
		logger.Logger.Errorw("error inserting order", "error", "no rows inserted")

		return fmt.Errorf("no rows inserted")
	}

	return nil
}

func (r *Repository) GetOrder(ctx context.Context, order int) (userID int, err error) {
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

func (r *Repository) UpdateStatus(ctx context.Context, order int, status string) error {
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
