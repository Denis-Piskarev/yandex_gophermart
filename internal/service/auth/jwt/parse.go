package jwt

import (
	"errors"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
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
		j.logger.Errorw("error while parsing token", "error", err)

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
		j.logger.Errorw("variable type conversion error")

		return 0, fmt.Errorf("variable type conversion error")
	}

	return int(userId), nil
}
