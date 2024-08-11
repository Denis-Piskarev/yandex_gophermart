package internal

// TokenInterface - operating with tokens
type TokenInterface interface {
	// GenerateToken - generating tokens token by userId. Returning token and error
	GenerateToken(login string) (string, error)
	// CheckToken - checking token. Returning login and error
	CheckToken(token string) (string, error)
}
