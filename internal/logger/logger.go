package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func NewLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.StacktraceKey = "" // disabling tracing error

	newDev, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := newDev.Sync(); err != nil {
			log.Println(fmt.Errorf("failed to sync %w", err))
		}
	}()

	Logger = newDev.Sugar()
}
