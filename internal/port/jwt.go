package ports

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwtInterface interface {
	Sign(userID string, role string) (string, error)
	Parse(tokenStr string) (*Claims, error)
}
