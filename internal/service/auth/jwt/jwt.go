package jwt

import "go.uber.org/zap"

const SecretKey = "not secret"

type JWT struct {
	logger *zap.SugaredLogger
}

func NewJWT(logger *zap.SugaredLogger) *JWT {
	return &JWT{logger: logger}
}
