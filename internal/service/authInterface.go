package service

import (
	"github.com/DenisquaP/yandex_gophermart/internal/service/auth"
)

// AuthInterface - operate with user auth and registration
type AuthInterface interface {
	auth.TokenInterface
	// RegisterUser - register new user in system by login and password. Returns error
	RegisterUser(username, password string) error
	// GetHashedPassword - gets hash from password. Returns hashed password and error
	GetHashedPassword(password string) (string, error)
	// LoginUser - checks pair login-password. Returns error
	LoginUser(username, password string) error
}
