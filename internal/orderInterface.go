package internal

import (
	"context"

	"github.com/DenisquaP/yandex_gophermart/internal/models/orders"
)

// OrderInterface - operates with user`s orders
type OrderInterface interface {
	// UploadOrder - upload new order into service
	UploadOrder(ctx context.Context, userId string, order orders.Order) error
	// GetOrders - gets orders of current user. Returns slice of orders and error
	GetOrders(ctx context.Context, userId string) ([]orders.Order, error)
}
