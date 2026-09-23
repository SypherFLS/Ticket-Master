package auth

import (
	"time"
	"tmaster/internal/apperrors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret []byte) *JWTManager {
	return &JWTManager{
		secret: secret,
	}
}

type claims struct {
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`

	jwt.RegisteredClaims
}

func (j *JWTManager) Generate(userID int, UserRole string) (string, error) {
	claims := claims{
		UserID:   userID,
		UserRole: UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}

func (j *JWTManager) Validate(tokenString string) (int, string, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, apperrors.WrongSignMethod
			}

			return j.secret, nil
		},
	)

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", apperrors.InvalidToken
	}

	return claims.UserID, claims.UserRole, nil
}
