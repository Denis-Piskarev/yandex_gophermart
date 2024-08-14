package jwt

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
)

// ParseToken - checking token. Returning login and error
func (j *JWT) ParseToken(tokenString string) (int, error) {
	// parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SecretKey), nil
	})
	if err != nil {
		logger.Logger.Errorw("error while parsing token", "error", err)

		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, models.NewCustomError("token is expired", http.StatusUnauthorized)
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, models.NewCustomError("invalid token", http.StatusBadRequest)
		}

		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}

	// get userId from claims
	userId, ok := claims["userId"].(float64)
	if !ok {
		logger.Logger.Errorw("variable type conversion error")

		return 0, fmt.Errorf("variable type conversion error")
	}

	return int(userId), nil
}
