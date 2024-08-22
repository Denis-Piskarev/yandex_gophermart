package user

import (
	"crypto/sha256"
	"encoding/hex"
)

// GetHashedPassword - generate hash from password to store secure data
func (u *User) GetHashedPassword(password string) string {
	hash := sha256.New()

	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil))
}
