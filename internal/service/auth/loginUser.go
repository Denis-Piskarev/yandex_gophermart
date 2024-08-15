package auth

import (
	"context"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customErrors"
	"net/http"
)

func (a *UserAuth) LoginUser(ctx context.Context, username, password string) (string, error) {
	// generating hash from password
	hashPassword := a.GetHashedPassword(password)

	userId, err := a.db.Login(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	// if userId == 0 => user not exists
	if userId == 0 {
		err := customErrors.CustomError{
			Err:        "user not found",
			StatusCode: http.StatusUnauthorized,
		}

		logger.Logger.Errorw("user not found", "username", username)

		return "", err
	}

	// generating token
	token, err := a.token.GenerateToken(userId)

	return token, err
}
