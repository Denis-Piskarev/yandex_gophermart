package middlewares

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"net/http"
	"strconv"
)

func IsAuthorized(services *internal.ServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("gopherToken")
			if err != nil {
				logger.Logger.Errorw("error getting token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			userId, err := services.ParseToken(token.Value)
			if err != nil || userId == 0 {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			r.Header.Set("userId", strconv.Itoa(userId))

			next.ServeHTTP(w, r)
		})
	}

}
