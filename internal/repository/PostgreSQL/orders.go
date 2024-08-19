package PostgreSQL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UploadOrder(ctx context.Context, userId int, order *modelsOrder.Order) error {
	return nil
}

func (r *Repository) GetOrder(ctx context.Context, order string) (userId int, err error) {
	query := `SELECT user_id FROM orders WHERE number = $1`

	if err := r.db.QueryRow(ctx, query, order).Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		logger.Logger.Errorw("error in getting orders", "userId", userId, "err", err)

		return 0, err
	}

	return userId, err
}

func (r *Repository) GetOrders(ctx context.Context, userId int) ([]*modelsOrder.Order, error) {
	query := `SELECT number, status, accural, uploaded_at FROM orders WHERE user_id = $1`
	var orders []*modelsOrder.Order

	// use null int because accural can be NULL
	var accural sql.NullInt64

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		logger.Logger.Errorw("error in getting orders", "userId", userId, "err", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order modelsOrder.Order
		if err := rows.Scan(&order.Number, &order.Status, &accural, &order.UploadedAt); err != nil {
			logger.Logger.Errorw("error in getting orders", "userId", userId, "err", err)

			return nil, err
		}

		// check accural for not NULL value
		if accural.Valid {
			order.Accural = int(accural.Int64)
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
