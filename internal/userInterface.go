package internal

import (
	"context"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

// AuthInterface - operate with user user and registration
type AuthInterface interface {
	// RegisterUser - register new user in system by login and password. Returns error and token
	RegisterUser(ctx context.Context, username, password string) (string, error)
	// GetHashedPassword - gets hash from password. Returns hashed password and error
	GetHashedPassword(password string) string
	// LoginUser - checks pair login-password. Returns token and error
	LoginUser(ctx context.Context, username, password string) (string, error)
	// GetWithdrawals - gets withdrawals of user. Returns custom error from models if user has no withdrawals
	GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error)
	// Withdraw - withdraws balance of user
	Withdraw(ctx context.Context, userID, sum int, order string) error
}
