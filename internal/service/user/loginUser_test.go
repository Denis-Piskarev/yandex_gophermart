package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	mocks_repo "github.com/DenisquaP/yandex_gophermart/internal/repository/mocks"
	"github.com/DenisquaP/yandex_gophermart/internal/service/jwt"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.NewLogger()
}

func TestUser_LoginUser(t *testing.T) {
	login := "login"
	password := "pass"
	userID := 1

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	hashedPassword := u.GetHashedPassword(password)

	db.EXPECT().Login(ctx, login, hashedPassword).Return(userID, nil)

	_, err := u.LoginUser(ctx, login, password)
	require.NoError(t, err)
}

func TestUser_LoginUser_NotExists(t *testing.T) {
	login := "login"
	password := "pass"
	userID := 0

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	hashedPassword := u.GetHashedPassword(password)

	db.EXPECT().Login(ctx, login, hashedPassword).Return(userID, nil)

	_, err := u.LoginUser(ctx, login, password)
	require.Error(t, err, "user not found")
}

func TestUser_LoginUser_DBError(t *testing.T) {
	login := "login"
	password := "pass"
	userID := 0

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db := mocks_repo.NewMockDBStore(ctrl)

	j := jwt.NewJWT()
	u := NewUserAuth(db, j)
	hashedPassword := u.GetHashedPassword(password)

	db.EXPECT().Login(ctx, login, hashedPassword).Return(userID, fmt.Errorf("err"))

	_, err := u.LoginUser(ctx, login, password)
	require.Error(t, err)
}
