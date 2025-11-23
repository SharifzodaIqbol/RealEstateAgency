package main

import (
	"example-app/pkg/estate"
	"example-app/pkg/store"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var tokenAuth *jwtauth.JWTAuth

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
	}
	tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)
}

func main() {
	initEnv()
	err := estate.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := sqlx.Connect("postgres", os.Getenv("CONNECT_SQL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Успешно подключено к базе данных.")
	defer db.Close()
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		jwtauth.Verifier(tokenAuth),
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	)
	auth := store.NewStoreDB(db)
	login := store.NewStore(auth, tokenAuth)
	r.Post("/register", auth.Register)
	r.Post("/login", login.Login)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator)
		r.Route("/properties", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/", estate.Create[estate.Property])
				r.Get("/", estate.Read[estate.Property])
				r.Get("/my", estate.GetMyData[estate.Property])
				r.Get("/{id}", estate.GetByID[estate.Property])
				r.Put("/{id}", estate.Update[estate.Property])
				r.Delete("/{id}", estate.Delete[estate.Property])
			})
		})
		r.Route("/purchases", func(r chi.Router) {
			r.Get("/", estate.Read[estate.Purchase])
			r.Get("/my", estate.GetMyData[estate.Purchase])
			r.Get("/{id}", estate.GetByID[estate.Purchase])
			r.Group(func(r chi.Router) {
				r.Use(estate.RequireAdminOrAgent)
				r.Post("/", estate.Create[estate.Purchase])
				r.Put("/{id}", estate.Update[estate.Purchase])
			})
		})
		r.Route("/sales", func(r chi.Router) {
			r.Get("/", estate.Read[estate.Sale])
			r.Get("/my", estate.GetMyData[estate.Sale])
			r.Get("/{id}", estate.GetByID[estate.Sale])
			r.Group(func(r chi.Router) {
				r.Use(estate.RequireAdminOrAgent)
				r.Post("/", estate.Create[estate.Sale])
				r.Put("/{id}", estate.Update[estate.Sale])
			})
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Authenticator)
		r.Use(estate.RequireAdmin)
		r.Route("/admin", func(r chi.Router) {
			// Управление пользователями
			r.Get("/users", estate.Read[estate.User])
			r.Put("/users/{id}/role", estate.Update[estate.User])
			r.Delete("/users/{id}", estate.Delete[estate.User])

			// Управление системой
			r.Delete("/purchases/{id}", estate.Delete[estate.Purchase])
			r.Delete("/sales/{id}", estate.Delete[estate.Sale])
		})
	})
	fmt.Println("Server started on :3000")
	http.ListenAndServe(":3000", r)
}
