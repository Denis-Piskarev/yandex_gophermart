package user

import (
	"context"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (u *User) GetWithdrawals(ctx context.Context, userId int) ([]*modelsUser.Withdrawals, error) {
	return u.db.GetWithdrawals(ctx, userId)
}
