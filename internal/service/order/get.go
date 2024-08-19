package order

import (
	"context"

	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
)

func (o *Order) GetOrders(ctx context.Context, userId int) ([]*modelsOrder.Order, error) {
	return o.db.GetOrders(ctx, userId)
}