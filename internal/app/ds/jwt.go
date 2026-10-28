package ds

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
}