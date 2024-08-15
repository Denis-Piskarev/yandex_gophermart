package PostgreSQL

import (
	"context"
	"database/sql"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
)

func (r *Repository) UploadOrder(ctx context.Context, userId int, order *modelsOrder.Order) error {
	return nil
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

func (r *Repository) UpdateStatus(ctx context.Context, userId int, order *modelsOrder.Order) error {
	return nil
}
