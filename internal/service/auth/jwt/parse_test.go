package jwt

import (
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestJWT_ParseToken(t *testing.T) {
	nd, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() { _ = nd.Sync() }()

	userId := 1

	jwt := NewJWT(nd.Sugar())

	token, err := jwt.GenerateToken(userId)
	require.NoError(t, err)

	require.NotEmpty(t, token)

	user, err := jwt.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, userId, user)
}

func TestJWT_ParseToken_expiredToken(t *testing.T) {
	nd, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() { _ = nd.Sync() }()

	jwt := NewJWT(nd.Sugar())

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjM1MzU5MjQsImlhdCI6MTcyMzUzOTUyNCwidXNlcklkIjoxfQ.8xN6IK8gbRXyroc_ufk-BqOQwGFyEomfdj_whX9Vy1g"

	_, err = jwt.ParseToken(expiredToken)
	var cErr models.CustomError
	if errors.As(err, &cErr) {
		require.Equal(t, http.StatusUnauthorized, cErr.StatusCode)
	}
}

func TestJWT_ParseToken_wrongToken(t *testing.T) {
	nd, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() { _ = nd.Sync() }()

	jwt := NewJWT(nd.Sugar())

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err = jwt.ParseToken(expiredToken)
	var cErr models.CustomError
	if errors.As(err, &cErr) {
		require.Equal(t, http.StatusBadRequest, cErr.StatusCode)
	}
}
