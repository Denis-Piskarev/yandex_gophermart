package userAuth

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserAuth_Register(t *testing.T) {
	a := UserAuth{}
	err := a.Register("a", "a")
	require.NoError(t, err)

}
