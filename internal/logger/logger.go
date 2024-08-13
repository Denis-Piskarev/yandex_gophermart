package logger

import (
	"fmt"
	"go.uber.org/zap"
	"log"
)

var Logger *zap.SugaredLogger

func NewLogger() {
	newDev, err := zap.NewDevelopment()
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
