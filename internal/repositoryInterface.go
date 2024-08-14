package internal

import "context"

// DBStore - interface of database
type DBStore interface {
	// Register - register new user in service
	Register(ctx context.Context, login, password string) error
	// CheckLogin - check login for exist. Returns error if login already occupied
	CheckLogin(ctx context.Context, login string) error
}
