package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DenisquaP/yandex_gophermart/internal/models"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	mocks_repo "github.com/DenisquaP/yandex_gophermart/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBalance_GetOrders_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := 1

	db := mocks_repo.NewMockDBStore(ctrl)
	db.EXPECT().GetOrders(ctx, userID).Return(nil, fmt.Errorf("err"))

	o := NewOrder(db, "")

	_, err := o.GetOrders(ctx, userID)
	require.Error(t, err)
}

func TestBalance_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := 1
	orders := []*modelsOrder.Order{{
		Number:     "s12312",
		Status:     models.PROCESSED,
		Accrual:    123.2,
		UploadedAt: time.Now(),
	}}

	db := mocks_repo.NewMockDBStore(ctrl)
	db.EXPECT().GetOrders(ctx, userID).Return(orders, fmt.Errorf("err"))

	o := NewOrder(db, "")

	ordersS, err := o.GetOrders(ctx, userID)
	require.Error(t, err)
	require.Equal(t, orders, ordersS)
}
