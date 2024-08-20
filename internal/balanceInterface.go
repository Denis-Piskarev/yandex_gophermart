package internal

import (
	"context"

	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"
)

// BalanceKeeper - operate with orders
type BalanceKeeper interface {
	// GetBalance - gets balance of user
	GetBalance(ctx context.Context, userID int) (modelsBalance.Balance, error)
	// WithdrawBalance -
	WithdrawBalance(ctx context.Context, userID int, amount float64) error
}
