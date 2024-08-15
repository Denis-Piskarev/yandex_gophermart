package jwt

import (
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customErrors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func init() {
	logger.NewLogger()
}

func TestJWT_ParseToken(t *testing.T) {
	userId := 1

	jwt := NewJWT()

	token, err := jwt.GenerateToken(userId)
	require.NoError(t, err)

	require.NotEmpty(t, token)

	user, err := jwt.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, userId, user)
}

func TestJWT_ParseToken_expiredToken(t *testing.T) {
	jwt := NewJWT()

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjM1MzU5MjQsImlhdCI6MTcyMzUzOTUyNCwidXNlcklkIjoxfQ.8xN6IK8gbRXyroc_ufk-BqOQwGFyEomfdj_whX9Vy1g"

	_, err := jwt.ParseToken(expiredToken)
	var cErr customErrors.CustomError
	if errors.As(err, &cErr) {
		require.Equal(t, http.StatusUnauthorized, cErr.StatusCode)
	}
}

func TestJWT_ParseToken_wrongToken(t *testing.T) {
	jwt := NewJWT()

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := jwt.ParseToken(expiredToken)
	var cErr customErrors.CustomError
	if errors.As(err, &cErr) {
		require.Equal(t, http.StatusBadRequest, cErr.StatusCode)
	}
}
