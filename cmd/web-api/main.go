package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// Ваши импорты Auth Service
	auth_model "auth-service/internal/auth/model" // ПРИМЕЧАНИЕ: Замените 'auth-service' на 'course-project'
	auth_service "auth-service/internal/auth/service"

	// НОВЫЕ ИМПОРТЫ для Недвижимости
	estate_handler "auth-service/internal/estate/handler"
	estate_repository "auth-service/internal/estate/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

type contextKey string

const userContextKey contextKey = "agent"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env. Используются системные переменные окружения.")
	}

	const serverPort = ":8081"
	const tokenExp = time.Minute * 15
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env. Используются системные переменные окружения.")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Переменная окружения JWT_SECRET не установлена!")
	}
	jwtService := auth_service.NewJWTService(jwtSecret, tokenExp, time.Hour*24*7)
	dbConnString := fmt.Sprintf("user=postgres password=%s dbname=auth sslmode=disable", os.Getenv("mypass"))
	if dbConnString == "" {
		log.Fatal("Переменная окружения WEBAPI_DB_CONNECTION_STRING не установлена!")
	}

	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД Web-API: %v", err)
	}
	defer db.Close()
	log.Println("Успешное подключение к базе данных Web API.")

	estateRepo := estate_repository.NewSQLEstateRepository(db)
	estateHandler := estate_handler.NewEstateHandler(estateRepo)

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

	r.Route("/api/v1/estate", func(r chi.Router) {
		// Все роуты недвижимости требуют авторизации
		r.Use(AuthMiddleware(jwtService))

		// CRUD для Объектов Недвижимости (требует, например, роли 'agent')
		r.Post("/properties", estateHandler.CreateProperty)
		r.Get("/properties/{id}", estateHandler.GetPropertyByID)
		r.Get("/properties", estateHandler.ListProperties)

		r.Post("/sales", estateHandler.CreateSale)
		r.Post("/purchases", estateHandler.CreatePurchase)
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

			// Детальная проверка прав
			if strings.HasPrefix(r.URL.Path, "/api/v1/estate") {
				if claims.RoleName != "agent" && r.Method == "POST" {
					http.Error(w, "Forbidden: agent role required", http.StatusForbidden)
					return
				}
			}

			log.Printf("User authenticated: ID=%d, Role=%s", claims.UserID, claims.RoleName)

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleSecretData(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(userContextKey).(*auth_model.CustomClaims)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Secret data for User ID: " + strconv.Itoa(claims.UserID)))
	w.Write([]byte(" | Role: " + claims.RoleName))
}
