package service

import "auth-service/internal/auth/model"

type TokenService interface {
	GenerateTokens(userID int, RoleName string) (accessToken string, refreshToken string, err error)
	ValidateToken(tokenString string) (*model.CustomClaims, error)
}
