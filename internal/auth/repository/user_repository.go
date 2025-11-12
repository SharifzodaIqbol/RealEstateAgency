package repository

import "auth-service/internal/auth/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	FindRoleByName(roleName string) (int, error)
	FindUserByIdentifier(identifier string) (*model.User, string, error)
	GetRoleIDByName(roleName string) (int, error)
}
