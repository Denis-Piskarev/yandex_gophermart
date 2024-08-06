package orders

import (
	"context"

	"github.com/DenisquaP/yandex_gophermart/internal/models/orders"
)

// Manager - operates with user`s orders
type Manager interface {
	// Upload - upload new order into service
	Upload(ctx context.Context, userId string, order orders.Order) error
	// Get - gets orders of current user. Returns slice of orders and error
	Get(ctx context.Context, userId string) ([]orders.Order, error)
}
