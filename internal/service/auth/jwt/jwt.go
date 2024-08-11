package jwt

import "go.uber.org/zap"

type JWT struct {
	logger *zap.SugaredLogger
}

func NewJWT(logger *zap.SugaredLogger) *JWT {
	return &JWT{logger: logger}
}
