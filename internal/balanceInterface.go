package internal

import "context"

// BalanceKeeper - operate with orders
type BalanceKeeper interface {
	// GetBalance - gets balance of user
	GetBalance(ctx context.Context, userId string) (float64, error)
	// WithdrawBalance -
	WithdrawBalance(ctx context.Context, userId string, amount float64) error
}
