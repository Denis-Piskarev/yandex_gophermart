package jwt

const SecretKey = "not secret"

type JWT struct {
}

func NewJWT() *JWT {
	return &JWT{}
}
