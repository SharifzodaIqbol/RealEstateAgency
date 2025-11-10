package model

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}
