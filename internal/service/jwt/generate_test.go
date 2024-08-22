package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJWT_GenerateToken(t *testing.T) {
	jwt := NewJWT()

	token, err := jwt.GenerateToken(1)
	require.NoError(t, err)
	fmt.Println(token)

	require.NotEmpty(t, token)
}
