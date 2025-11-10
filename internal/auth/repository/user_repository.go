package repository

import "auth-service/internal/auth/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, string, error)
	FindRoleByName(roleName string) (int, error)
}
