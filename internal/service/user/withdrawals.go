package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (u *User) GetWithdrawals(ctx context.Context, userID int) ([]*modelsUser.Withdrawals, error) {
	return u.db.GetWithdrawals(ctx, userID)
}

func (u *User) Withdraw(ctx context.Context, userID int, sum float32, order string) error {
	orderInt, err := strconv.Atoi(order)
	if err != nil {
		logger.Logger.Errorw("Error converting order to int", "order", order)
		cErr := customerrors.NewCustomError("unsupported type of order", http.StatusUnprocessableEntity)

		return cErr
	}

	userIDOrder, err := u.db.GetOrder(ctx, orderInt)
	if err != nil {
		return err
	}

	//if userIDOrder != userID {
	//	logger.Logger.Errorw("user id not match")
	//	cErr := customerrors.NewCustomError("user id not match", http.StatusUnprocessableEntity)
	//
	//	return cErr
	//}

	return u.db.Withdraw(ctx, userIDOrder, sum, orderInt)
}
