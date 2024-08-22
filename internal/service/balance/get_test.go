package balance

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	modelsBalance "github.com/DenisquaP/yandex_gophermart/internal/models/balance"
	mocks_repo "github.com/DenisquaP/yandex_gophermart/internal/repository/mocks"

	"github.com/golang/mock/gomock"
)

func TestBalance_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := 1
	balance := modelsBalance.Balance{
		Current:   1.2,
		Withdrawn: 30.34,
	}

	db := mocks_repo.NewMockDBStore(ctrl)
	db.EXPECT().GetBalance(ctx, userID).Return(balance, nil)

	b := NewBalance(db)

	userBalance, err := b.GetBalance(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, balance, userBalance)
}

func TestBalance_GetBalance_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := 1
	balance := modelsBalance.Balance{}

	db := mocks_repo.NewMockDBStore(ctrl)
	db.EXPECT().GetBalance(ctx, userID).Return(balance, fmt.Errorf("err"))

	b := NewBalance(db)

	_, err := b.GetBalance(ctx, userID)
	require.Error(t, err)
}
