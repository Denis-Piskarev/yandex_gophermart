package auth

import (
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth/tokens"
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth/userAuth"
)

type Manager interface {
	tokens.Manager
	userAuth.Manager
}
