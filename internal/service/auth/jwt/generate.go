package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
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
		j.logger.Errorw("error generating token", "err", err)

		return "", err
	}

	return tokenString, nil
}
