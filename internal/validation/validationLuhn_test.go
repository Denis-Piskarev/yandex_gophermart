package validation

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidateLuhn_Ok(t *testing.T) {
	order := "12345678903"
	valid := ValidateLuhn(order)
	require.True(t, valid)
}

func Test_ValidateLuhn_NotValid(t *testing.T) {
	order := "1234567893"
	valid := ValidateLuhn(order)
	require.False(t, valid)
}
