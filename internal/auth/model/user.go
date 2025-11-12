package model

import "time"

type User struct {
	ID           int       `db:"id"`
	UserName     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	RoleID       int       `db:"role_id"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
