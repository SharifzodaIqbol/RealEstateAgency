package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	auth_model "auth-service/internal/auth/model"
	auth_service "auth-service/internal/auth/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type contextKey string

const userContextKey contextKey = "user"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env. Используются системные переменные окружения.")
	}

	const serverPort = ":8081"
	const tokenExp = time.Minute * 15

	jwtSecret := os.Getenv("jwtSecret")
	if jwtSecret == "" {
		log.Fatal("Переменная окружения jwtSecret не установлена!")
	}
	jwtService := auth_service.NewJWTService(jwtSecret, tokenExp, time.Hour*24*7)

	// 3. Настройка Роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Web API is running"))
	})

	r.Route("/api/v1/data", func(r chi.Router) {
		r.Use(AuthMiddleware(jwtService))

		r.Get("/secret", handleSecretData)
	})

	log.Printf("Web API запущен на порту %s", serverPort)
	if err := http.ListenAndServe(serverPort, r); err != nil {
		log.Fatalf("Ошибка запуска Web API: %v", err)
	}
}
func AuthMiddleware(ts auth_service.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			claims, err := ts.ValidateToken(tokenString)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func handleSecretData(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userContextKey).(*auth_model.CustomClaims)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Secret data for User ID: " + strconv.Itoa(claims.UserID)))
	w.Write([]byte(" | Role: " + claims.RoleName))
}
