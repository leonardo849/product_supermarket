package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	AuthID         string `json:"id"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}