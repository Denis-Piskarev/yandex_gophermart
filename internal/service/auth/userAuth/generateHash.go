package userAuth

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetHashedPassword - generate hash from password to store secure data
func (a *UserAuth) GetHashedPassword(password string) (string, error) {
	hash := sha256.New()

	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil)), nil
}
