package internal

// DBStore - interface of database
type DBStore interface {
	// Register - register new user in service
	Register(login, password string) error
}
