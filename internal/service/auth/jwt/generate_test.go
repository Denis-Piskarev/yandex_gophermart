package jwt

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJWT_GenerateToken(t *testing.T) {
	jwt := NewJWT()

	token, err := jwt.GenerateToken(1)
	require.NoError(t, err)
	fmt.Println(token)

	require.NotEmpty(t, token)
}
