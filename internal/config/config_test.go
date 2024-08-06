package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigWithEnvValues(t *testing.T) {
	// Setting env values
	t.Setenv("RUN_ADDRESS", "localhost:12323")
	t.Setenv("DATABASE_URI", "postgres")
	t.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:2222")

	// Getting config
	config, err := NewConfig()
	require.NoError(t, err)

	// Checking values
	require.Equal(t, "localhost:12323", config.RunAddress)
	require.Equal(t, "postgres", config.DatabaseUri)
	require.Equal(t, "localhost:2222", config.AccuralSystemAddress)
}

func TestNewConfigDefault(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = append([]string{os.Args[0], "-a=localhost:8080", "-d=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "-r="}, os.Args...)

	// Getting config
	config, err := NewConfig()
	require.NoError(t, err)

	// Checking values
	require.Equal(t, "localhost:8080", config.RunAddress)
	require.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", config.DatabaseUri)
	require.Equal(t, "", config.AccuralSystemAddress)
}
