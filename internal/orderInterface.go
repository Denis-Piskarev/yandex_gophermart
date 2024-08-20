package internal

import (
	"context"
	"github.com/DenisquaP/yandex_gophermart/internal/models/orders"
)

// OrderInterface - operates with user`s orders
type OrderInterface interface {
	// UploadOrder - upload new order into service, returns status code if error is nil
	UploadOrder(ctx context.Context, userID int, order string) (int, error)
	// GetOrders - gets orders of current user. Returns slice of orders and error
	GetOrders(ctx context.Context, userID int) ([]*orders.Order, error)
}
