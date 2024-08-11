package userAuth

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestUserAuth_Register(t *testing.T) {
	nd, err := zap.NewDevelopment()
	require.NoError(t, err)

	logger := nd.Sugar()

	a := UserAuth{logger: logger}
	err = a.Register("a", "a")
	require.NoError(t, err)

}
