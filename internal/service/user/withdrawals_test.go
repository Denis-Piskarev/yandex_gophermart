package user

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	userModels "github.com/DenisquaP/yandex_gophermart/internal/models/users"
	mocks_repo "github.com/DenisquaP/yandex_gophermart/internal/repository/mocks"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUser_GetWithdrawals(t *testing.T) {
	userID := 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	w := []*userModels.Withdrawals{
		{},
	}

	db.EXPECT().GetWithdrawals(ctx, userID).Return(w, nil)

	withdrawals, err := u.GetWithdrawals(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, withdrawals, w)
}

func TestUser_GetWithdrawals_NoWithdrawals(t *testing.T) {
	userID := 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	cErr := customerrors.NewCustomError("user has no withdrawals", http.StatusNoContent)

	db.EXPECT().GetWithdrawals(ctx, userID).Return(nil, cErr)

	withdrawals, err := u.GetWithdrawals(ctx, userID)
	require.ErrorIs(t, err, cErr)
	require.Nil(t, withdrawals)
}

func TestUser_GetWithdrawals_DBError(t *testing.T) {
	userID := 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	db.EXPECT().GetWithdrawals(ctx, userID).Return(nil, fmt.Errorf("err"))

	withdrawals, err := u.GetWithdrawals(ctx, userID)
	require.Error(t, err)
	require.Nil(t, withdrawals)
}

func TestUser_Withdraw(t *testing.T) {
	userID := 1
	var sum float32 = 123.2
	order := "123213"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	db.EXPECT().Withdraw(ctx, userID, sum, order).Return(nil)

	err := u.Withdraw(ctx, userID, sum, order)
	require.NoError(t, err)
}

func TestUser_Withdraw_NotEnoughBalance(t *testing.T) {
	userID := 1
	var sum float32 = 123.2
	order := "123213"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	cErr := customerrors.NewCustomError("not enough balance", http.StatusPaymentRequired)
	db.EXPECT().Withdraw(ctx, userID, sum, order).Return(cErr)

	err := u.Withdraw(ctx, userID, sum, order)
	require.ErrorIs(t, err, cErr)
}

func TestUser_Withdraw_ErrorInDB(t *testing.T) {
	userID := 1
	var sum float32 = 123.2
	order := "123213"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)

	db.EXPECT().Withdraw(ctx, userID, sum, order).Return(fmt.Errorf("err"))

	err := u.Withdraw(ctx, userID, sum, order)
	require.Error(t, err)
}
