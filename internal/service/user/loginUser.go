package user

import (
	"context"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
)

func (u *User) LoginUser(ctx context.Context, username, password string) (string, error) {
	// generating hash from password
	hashPassword := u.GetHashedPassword(password)

	userID, err := u.db.Login(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	// if userID == 0 => user not exists
	if userID == 0 {
		err := customerrors.CustomError{
			Err:        "user not found",
			StatusCode: http.StatusUnauthorized,
		}

		logger.Logger.Errorw("user not found", "username", username)

		return "", err
	}

	// generating token
	return u.token.GenerateToken(userID)
}
