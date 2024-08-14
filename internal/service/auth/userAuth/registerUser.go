package userAuth

import (
	"context"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
)

func (a *UserAuth) RegisterUser(ctx context.Context, username string, password string) error {
	// check if login already occupied
	if err := a.db.CheckLogin(ctx, username); err != nil {
		return err
	}

	// generating hash from password
	hashPassword, err := a.GetHashedPassword(password)
	if err != nil {
		logger.Logger.Errorw("error while generating hash",
			"error", err)

		return err
	}

	// registering user
	return a.db.Register(ctx, username, hashPassword)
}
