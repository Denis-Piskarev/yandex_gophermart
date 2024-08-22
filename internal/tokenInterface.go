package internal

// TokenInterface - operating with tokens
type TokenInterface interface {
	// GenerateToken - generating tokens token by userID. Returning token and error
	GenerateToken(userID int) (string, error)
	// ParseToken - parsing token. Returning userID and error
	ParseToken(token string) (int, error)
}
