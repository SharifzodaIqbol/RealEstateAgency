// Package store
package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

func (s *StoreDB) Register(w http.ResponseWriter, r *http.Request) {
	regReq := RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println(regReq)
	if regReq.UserName == "" || regReq.Email == "" || regReq.Password == "" {
		http.Error(w, "Заполните все поля!", http.StatusBadRequest)
		return
	}

	var email string
	s.db.Get(&email, "SELECT email FROM users WHERE email = $1", regReq.Email)
	if email != "" {
		http.Error(w, "Пользователь с этим email существует!", http.StatusConflict)
		return
	}
	hashedPassword, err := HashPassword(regReq.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = s.db.Exec(
		"INSERT INTO users (username, email, password_hash, role_id) VALUES ($1, $2, $3, $4)",
		regReq.UserName, regReq.Email, hashedPassword, 3,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Вы успешно зарегистрировались!",
	})
}

func (s *Store) Login(w http.ResponseWriter, r *http.Request) {
	loginReq := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user := User{}
	s.storeDB.db.Get(
		&user,
		"SELECT id, username, email, password_hash, role_id FROM users WHERE email = $1",
		loginReq.Email,
	)
	fmt.Println(user)
	if user.Email == "" || !CheckPassword(user.Password, loginReq.Password) {
		http.Error(w, "Неправильно указан email или пароль", http.StatusUnauthorized)
		return
	}

	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"sub":     user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: tokenString,
		RoleID:      user.RoleID,
	})
}
func GetIDUser(r *http.Request) (int, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token: missing user_id")
	}
	return int(id), nil
}
