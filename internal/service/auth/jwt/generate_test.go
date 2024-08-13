package jwt

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestJWT_GenerateToken(t *testing.T) {
	nd, err := zap.NewDevelopment()
	require.NoError(t, err)
	defer func() { _ = nd.Sync() }()

	jwt := NewJWT(nd.Sugar())

	token, err := jwt.GenerateToken(1)
	require.NoError(t, err)
	fmt.Println(token)

	require.NotEmpty(t, token)
}
