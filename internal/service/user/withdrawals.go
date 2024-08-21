package user

import (
	"context"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (u *User) GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error) {
	return u.db.GetWithdrawals(ctx, userID)
}

func (u *User) Withdraw(ctx context.Context, userID int, sum float32, order string) error {
	userIDOrder, err := u.db.GetOrder(ctx, order)
	if err != nil {
		return err
	}

	//if userIDOrder != userID {
	//	logger.Logger.Errorw("user id not match")
	//	cErr := customerrors.NewCustomError("user id not match", http.StatusUnprocessableEntity)
	//
	//	return cErr
	//}

	return u.db.Withdraw(ctx, userIDOrder, sum, order)
}
