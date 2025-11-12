package service

import (
	"auth-service/internal/auth/model"
	"auth-service/internal/auth/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(username, email, password string) error
	LoginUser(identifier, password string) (accessToken, refreshToken string, err error)
}
type authService struct {
	userRepo     repository.UserRepository
	tokenService TokenService
}

func NewAuthService(userRepo repository.UserRepository, tokenService TokenService) *authService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *authService) RegisterUser(username, email, password string) error {
	defaultRoleID, err := s.userRepo.FindRoleByName("user")
	if err != nil {
		return errors.New("default role not found")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	newUser := &model.User{
		UserName:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       defaultRoleID,
		IsActive:     true,
	}
	return s.userRepo.CreateUser(newUser)
}
func (s *authService) LoginUser(identifier, password string) (accessToken, refreshToken string, err error) {

	user, roleName, err := s.userRepo.FindUserByIdentifier(identifier)

	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err = s.tokenService.GenerateTokens(user.ID, roleName)

	return accessToken, refreshToken, nil
}
