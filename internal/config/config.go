package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/caarlos0/env"
)

// Config - config of app
type Config struct {
	// RunAddress - address of server
	RunAddress string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	// DatabaseURI - uri of database
	DatabaseURI string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	// AccuralSystemAddress - address of accrual system
	AccuralSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
}

// NewConfig - returns config
func NewConfig() (*Config, error) {
	var cfg Config

	// Setting values by flags, if env not empty, using env
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "database address")
	flag.StringVar(&cfg.AccuralSystemAddress, "r", "", "accrual system address")

	if err := env.Parse(&cfg); err != nil {
		return &Config{}, err
	}

	flag.Parse()

	// if DatabaseDsn not empty, using it
	addr := strings.Split(cfg.DatabaseURI, " ")
	if len(addr) > 1 {
		cfg.DatabaseURI = getConnectionString(addr)
	}

	return &cfg, nil
}

// getConnectionString - returns postgres connection string
func getConnectionString(addr []string) string {
	var user string
	var pass string
	var host string
	var port string
	var dbname string

	for _, i := range addr {
		i = strings.Trim(i, `"`)
		variables := strings.Split(i, "=")
		switch variables[0] {
		case "user":
			user = variables[1]
		case "password":
			pass = variables[1]
		case "host":
			host = variables[1]
		case "port":
			port = variables[1]
		case "dbname":
			dbname = variables[1]
		}

	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
}
