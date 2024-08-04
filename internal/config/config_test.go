package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigWithEnvValues(t *testing.T) {
	// Setting env values
	err := os.Setenv("RUN_ADDRESS", "localhost:12323")
	require.NoError(t, err)
	err = os.Setenv("DATABASE_URI", "postgres")
	require.NoError(t, err)
	err = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:2222")
	require.NoError(t, err)

	// Getting config
	config, err := NewConfig()
	require.NoError(t, err)

	// Checking values
	require.Equal(t, "localhost:12323", config.RunAddress)
	require.Equal(t, "postgres", config.DatabaseUri)
	require.Equal(t, "localhost:2222", config.AccuralSystemAddress)
}

func TestNewConfigDefault(t *testing.T) {
	// Getting config
	config, err := NewConfig()
	require.NoError(t, err)

	// Checking values
	require.Equal(t, "localhost:8080", config.RunAddress)
	require.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", config.DatabaseUri)
	require.Equal(t, "", config.AccuralSystemAddress)
}
