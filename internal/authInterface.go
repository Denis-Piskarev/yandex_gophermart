package internal

import (
	"context"
)

// AuthInterface - operate with user auth and registration
type AuthInterface interface {
	// RegisterUser - register new user in system by login and password. Returns error
	RegisterUser(ctx context.Context, username, password string) error
	// GetHashedPassword - gets hash from password. Returns hashed password and error
	GetHashedPassword(password string) (string, error)
	// LoginUser - checks pair login-password. Returns error
	LoginUser(ctx context.Context, username, password string) error
}
