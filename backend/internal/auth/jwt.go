package auth

type JWTManager struct {
	Secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager {
		Secret: secret,
	}
}