package jwt

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"

	"github.com/golang-jwt/jwt/v5"
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
			return 0, customerrors.NewCustomError("token is expired", http.StatusUnauthorized)
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, customerrors.NewCustomError("invalid token", http.StatusBadRequest)
		}

		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}

	// get userID from claims
	userID, ok := claims["userID"].(float64)
	if !ok {
		logger.Logger.Errorw("variable type conversion error")

		return 0, fmt.Errorf("variable type conversion error")
	}

	return int(userID), nil
}
