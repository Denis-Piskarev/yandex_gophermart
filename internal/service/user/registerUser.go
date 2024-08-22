package user

import (
	"context"
)

func (u *User) RegisterUser(ctx context.Context, username string, password string) (string, error) {
	// check if login already occupied
	if err := u.db.CheckLogin(ctx, username); err != nil {
		return "", err
	}

	// generating hash from password
	hashPassword := u.GetHashedPassword(password)

	// registering user
	id, err := u.db.Register(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	return u.token.GenerateToken(id)
}
