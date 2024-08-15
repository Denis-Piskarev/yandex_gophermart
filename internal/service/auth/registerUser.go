package auth

import (
	"context"
)

func (a *UserAuth) RegisterUser(ctx context.Context, username string, password string) (string, error) {
	// check if login already occupied
	if err := a.db.CheckLogin(ctx, username); err != nil {
		return "", err
	}

	// generating hash from password
	hashPassword := a.GetHashedPassword(password)

	// registering user
	id, err := a.db.Register(ctx, username, hashPassword)
	if err != nil {
		return "", err
	}

	token, err := a.token.GenerateToken(id)

	return token, err
}
