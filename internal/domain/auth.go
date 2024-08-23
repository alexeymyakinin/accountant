package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserId    int
	UserEmail string
	jwt.RegisteredClaims
}
