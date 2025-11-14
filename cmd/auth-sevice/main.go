package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auth-service/internal/auth/handler"
	"auth-service/internal/auth/repository"
	"auth-service/internal/auth/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const (
		serverPort     = ":8080"
		accessExpMin   = time.Minute * 15
		refreshExpDays = time.Hour * 24 * 7
	)
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env. Используются системные переменные окружения.")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	dbConnString := fmt.Sprintf("user=postgres password=%s dbname=auth sslmode=disable", os.Getenv("mypass"))

	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	log.Println("Успешное подключение к базе данных!")

	userRepo := repository.NewSQLUserRepository(db)

	tokenService := service.NewJWTService(jwtSecret, accessExpMin, refreshExpDays)
	authService := service.NewAuthService(userRepo, tokenService)

	authHandler := handler.NewAuthHandler(authService)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	log.Printf("Сервер запущен на порту %s", serverPort)
	if err := http.ListenAndServe(serverPort, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
