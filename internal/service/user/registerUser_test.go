package user

import (
	"context"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	mocks_repo "github.com/DenisquaP/yandex_gophermart/internal/repository/mocks"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUserAuth_Register(t *testing.T) {
	login := "login"
	password := "pass"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	hashedPassword := u.GetHashedPassword(password)

	db.EXPECT().CheckLogin(ctx, login).Return(nil)
	db.EXPECT().Register(ctx, login, hashedPassword).Return(1, nil)

	token, err := u.RegisterUser(ctx, login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestUserAuth_Register_LoginOccupied(t *testing.T) {
	login := "login"
	password := "pass"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	cErr := customerrors.NewCustomError("user already exists", http.StatusConflict)

	db.EXPECT().CheckLogin(ctx, login).Return(cErr)

	token, err := u.RegisterUser(ctx, login, password)
	require.ErrorIs(t, err, cErr)
	require.Empty(t, token)
}

func TestUserAuth_Register_ErrorRegister(t *testing.T) {
	login := "login"
	password := "pass"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	hashedPassword := u.GetHashedPassword(password)

	db.EXPECT().CheckLogin(ctx, login).Return(nil)
	db.EXPECT().Register(ctx, login, hashedPassword).Return(0, fmt.Errorf("err"))

	token, err := u.RegisterUser(ctx, login, password)
	require.Error(t, err)
	require.Empty(t, token)
}
