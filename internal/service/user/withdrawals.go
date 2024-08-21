package user

import (
	"context"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (u *User) GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error) {
	return u.db.GetWithdrawals(ctx, userID)
}

func (u *User) Withdraw(ctx context.Context, userID int, sum float32, order string) error {
	return u.db.Withdraw(ctx, userID, sum, order)
}
