package endpoints

import (
	"net/http"
	"strconv"

	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
)

type Endpoints struct {
	services *internal.Service
}

func NewEndpoints(services *internal.Service) *Endpoints {
	return &Endpoints{
		services: services,
	}
}

func getUserIDFromHeader(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("userID")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Logger.Errorw("error while converting user id to int", "err", err)

		return 0, err
	}

	return userID, nil
}
