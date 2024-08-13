package internal

// TokenInterface - operating with tokens
type TokenInterface interface {
	// GenerateToken - generating tokens token by userId. Returning token and error
	GenerateToken(userId int) (string, error)
	// ParseToken - parsing token. Returning userId and error
	ParseToken(token string) (int, error)
}
