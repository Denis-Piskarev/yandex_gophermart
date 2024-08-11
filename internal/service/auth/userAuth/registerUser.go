package userAuth

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

func (a *UserAuth) Register(ctx context.Context, username string, password string) error {
	_, err := generateHash(password)
	if err != nil {
		a.logger.Errorw("error while generating hash",
			"error", err,
			zap.Int("status_code", http.StatusInternalServerError))

		return err
	}

	return nil
}
