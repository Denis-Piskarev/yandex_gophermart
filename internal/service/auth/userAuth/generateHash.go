package userAuth

import (
	"crypto/sha256"
	"encoding/hex"
)

// Func generate hash from password to store secure data
func generateHash(password string) (string, error) {
	hash := sha256.New()

	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil)), nil
}
