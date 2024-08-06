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
	defer func() {
		if err := newDev.Sync(); err != nil {
			log.Println(fmt.Errorf("failed to sync %w", err))
		}
	}()

	logger := newDev.Sugar()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalw("Failed to get config", "error", err)
	}

	runServer(cfg, logger)
}

// runServer - starts server
func runServer(cfg *config.Config, logger *zap.SugaredLogger) {
	logger.Infow(fmt.Sprintf("Starting server on %s...", cfg))
}
