package service

import (
	"auth-service/internal/auth/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey  []byte
	accessExp  time.Duration
	refreshExp time.Duration
}

func NewJWTService(secretKey string, accessExp time.Duration, refreshExp time.Duration) *JWTService {
	return &JWTService{
		secretKey:  []byte(secretKey),
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}
func (s *JWTService) GenerateTokens(userID int, roleName string) (string, string, error) {
	accessClaims := &model.CustomClaims{
		UserID:   userID,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return "", "", err
	}
	refreshClaims := &model.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return "", "", nil
	}
	return accessTokenString, refreshTokenString, nil
}
func (s *JWTService) ValidateToken(tokenString string) (*model.CustomClaims, error) {
	claims := &model.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	validatedClaims, ok := token.Claims.(*model.CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return validatedClaims, nil
}
