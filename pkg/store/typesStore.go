package store

import (
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type StoreDB struct {
	db *sqlx.DB
}
type Store struct {
	storeDB   *StoreDB
	tokenAuth *jwtauth.JWTAuth
}
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	RoleID    int       `json:"role_id" db:"role_id"`
}
type RegisterRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RoleID      int    `json:"role_id"`
}
type ProfileResponse struct {
	Name  string `json:"username"`
	Email string `json:"email"`
}

func NewStoreDB(db *sqlx.DB) *StoreDB {
	return &StoreDB{db: db}
}
func NewStore(storeDB *StoreDB, tokenAuth *jwtauth.JWTAuth) *Store {
	return &Store{storeDB: storeDB, tokenAuth: tokenAuth}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
