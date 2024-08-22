package internal

//go:generate mockgen -source=repositoryInterface.go -destination=repository/mocks/mocks_repo.go -package=mocks_repo

import (
	"context"

	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"
	"github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

// DBStore - interface of database
type DBStore interface {
	// Register - register new user in service
	Register(ctx context.Context, login, password string) (int, error)
	// CheckLogin - check login for exist. Returns error if login already occupied
	CheckLogin(ctx context.Context, login string) error
	// Login - login registered user
	Login(ctx context.Context, login, password string) (int, error)
	// UploadOrder - upload new order into service
	UploadOrder(ctx context.Context, userID int, order *orders.OrderAccrual) error
	// GetOrders - gets orders of current user. Returns slice of orders and error
	GetOrders(ctx context.Context, userID int) ([]*orders.Order, error)
	// UpdateStatus - updates status of order
	UpdateStatus(ctx context.Context, order string, status string) error
	// GetBalance - gets balance of current user. If user has no balance - return {0, 0}
	GetBalance(ctx context.Context, userID int) (modelsBalance.Balance, error)
	// GetWithdrawals - gets withdrawals of user. Returns custom error from models if user has no withdrawals
	GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error)
	// GetUserIdByOrder - gets userID and order if order exists
	GetUserIdByOrder(ctx context.Context, order string) (userID int, err error)
	// Withdraw - withdraws balance of user
	Withdraw(ctx context.Context, userID int, sum float32, order string) error
}
