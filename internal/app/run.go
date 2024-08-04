package app

import (
	"fmt"
	"log"

	"github.com/DenisquaP/yandex_gophermart/internal/config"
	"go.uber.org/zap"
)

// Run - starts app
func Run() {
	newDev, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer newDev.Sync()

	logger := newDev.Sugar()

	config, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("Failed to get config", "error", err)
	}

	runServer(config, logger)
}

// runServer - starts server
func runServer(config *config.Config, logger *zap.SugaredLogger) {
	logger.Infow(fmt.Sprintf("Starting server on %s...", config))
}
