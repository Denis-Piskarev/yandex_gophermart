package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateHash(t *testing.T) {
	a := User{}
	hash := a.GetHashedPassword("userPassword")

	require.NotEmpty(t, hash)
	fmt.Println(hash)
}
