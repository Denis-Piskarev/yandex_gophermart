package userAuth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateHash(t *testing.T) {
	a := UserAuth{}
	hash, err := a.GetHashedPassword("userPassword")
	require.NoError(t, err)

	require.NotEmpty(t, hash)
	fmt.Println(hash)
}
