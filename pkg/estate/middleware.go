package estate

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

func tableFromPath(path string) (string, error) {
	// нормализуем путь и берем первую часть: "/properties/123" -> "properties"
	clean := strings.Trim(path, "/")
	if clean == "" {
		return "", fmt.Errorf("empty path")
	}
	parts := strings.Split(clean, "/")
	table := parts[0]
	if !AllowedTables[table] {
		return "", fmt.Errorf("table not allowed: %s", table)
	}
	return table, nil
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		email, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var roleID int
		err := DB.Get(&roleID, "SELECT role_id FROM users WHERE email = $1", email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Println("RequireAdmin DB error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if roleID != 1 {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireAdminOrAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		email, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "User ID claim missing or invalid", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		var roleID int
		if err := DB.Get(&roleID, "SELECT role_id FROM users WHERE email = $1", email); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Println("RequireAdminOrAgent DB error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// если роль не admin(1) и не agent(2) — запрет
		if roleID != 1 && roleID != 2 {
			http.Error(w, "Agent access required", http.StatusForbidden)
			return
		}

		// если указан id в URL — проверим owner_id целевой записи
		idStr := chi.URLParam(r, "id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}
			table, err := tableFromPath(r.URL.Path)
			if err != nil {
				http.Error(w, "Invalid resource", http.StatusBadRequest)
				return
			}
			// безопасно составляем запрос — table уже whitelist'ирована
			query := fmt.Sprintf("SELECT owner_id FROM %s WHERE id = $1", table)
			var ownerID int
			if err := DB.Get(&ownerID, query, id); err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
				log.Println("RequireAdminOrAgent DB error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// если владелец есть и он не совпадает с текущим пользователем — запрет (за исключением админа)
			if ownerID != 0 && ownerID != userID && roleID != 1 {
				http.Error(w, "Недостаточно прав доступа", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
