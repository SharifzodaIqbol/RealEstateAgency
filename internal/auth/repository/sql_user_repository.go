package repository

import (
	"auth-service/internal/auth/model"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type SQLUserRepository struct {
	db *sqlx.DB
}

func NewSQLUserRepository(db *sqlx.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}
func (r *SQLUserRepository) FindRoleByName(roleName string) (int, error) {
	var roleID int
	query := `SELECT id FROM roles WHERE name = $1`

	err := r.db.Get(&roleID, query, roleName)

	if err == sql.ErrNoRows {
		return 0, errors.New("role not found")
	}
	return roleID, err
}

func (r *SQLUserRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users
              (username, email, password_hash, is_active, role_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.PasswordHash,
		user.IsActive,
		user.RoleID,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)

	return err
}
func (r *SQLUserRepository) FindUserByIdentifier(identifier string) (*model.User, string, error) {
	var user model.User
	var roleName string
	query := `SELECT 
				u.id, u.username, u.email, u.password_hash, u.role_id, u.is_active, u.created_at, u.updated_at, r.name AS role_name
              FROM users u
              JOIN roles r ON u.role_id = r.id
              WHERE u.username = $1 OR u.email = $1`

	row := r.db.QueryRowx(query, identifier)
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleName,
	)

	if err == sql.ErrNoRows {
		return nil, "", errors.New("user not found")
	}

	if err != nil {
		return nil, "", err
	}

	return &user, roleName, nil
}
func (r *SQLUserRepository) GetRoleIDByName(roleName string) (int, error) {
	var roleID int
	err := r.db.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)

	if err == sql.ErrNoRows {
		return 0, errors.New("role not found: " + roleName)
	}
	if err != nil {
		return 0, err
	}
	return roleID, nil
}
