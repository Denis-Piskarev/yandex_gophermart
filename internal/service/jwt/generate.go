package jwt

import (
	"time"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken - generating tokens token by userId. Returning token and error
func (j *JWT) GenerateToken(userId int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		logger.Logger.Errorw("error generating token", "err", err)

		return "", err
	}

	return tokenString, nil
}