package service

// BalanceKeeper - operate with orders
type BalanceKeeper interface {
	// GetBalance - gets balance of user
	GetBalance(userId string) (float64, error)
	// WithdrawBalance -
	WithdrawBalance(userId string, amount float64) error
}
