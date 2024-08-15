package endpoints

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"net/http"
	"strconv"
)

type Endpoints struct {
	services *internal.ServiceInterface
}

func NewEndpoints(services *internal.ServiceInterface) *Endpoints {
	return &Endpoints{
		services: services,
	}
}

func getUserIdFromHeader(r *http.Request) (int, error) {
	userIdStr := r.Header.Get("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logger.Logger.Errorw("error while converting user id to int", "err", err)

		return 0, err
	}

	return userId, nil
}
