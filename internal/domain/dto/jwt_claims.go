package dto

import "github.com/golang-jwt/jwt"

type Claims struct {
	UserID  uint
	IsAdmin bool
	jwt.StandardClaims
}
