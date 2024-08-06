package tokens

// Manager - operating with tokens
type Manager interface {
	// GenerateToken - generating tokens token by userId. Returning token and error
	GenerateToken(login string) (string, error)
	// CheckToken - checking token. Returning login and error
	CheckToken(token string) (string, error)
}
