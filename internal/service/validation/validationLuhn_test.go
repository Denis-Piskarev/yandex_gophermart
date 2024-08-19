package validation

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_IsValidLuhnNumber_Ok(t *testing.T) {
	order := "12345678903"
	valid := IsValidLuhnNumber(order)
	require.True(t, valid)
}

func Test_IsValidLuhnNumber_NotValid(t *testing.T) {
	order := "1234567893"
	valid := IsValidLuhnNumber(order)
	require.False(t, valid)
}
