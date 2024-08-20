package middlewares

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"net/http"
	"strconv"
)

func IsAuthorized(services *internal.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("gopherToken")
			if err != nil {
				logger.Logger.Errorw("error getting token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			userID, err := services.ParseToken(token.Value)
			if err != nil || userID == 0 {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			r.Header.Set("userID", strconv.Itoa(userID))

			next.ServeHTTP(w, r)
		})
	}

}
