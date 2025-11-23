// Package estate gfg
package estate

import (
	"database/sql"
	"encoding/json"
	"example-app/pkg/store"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

var AllowedTables = map[string]bool{
	"properties": true,
	"purchases":  true,
	"sales":      true,
	"users":      true,
}

func InitDB() error {
	var err error
	DB, err = sqlx.Connect("postgres", os.Getenv("CONNECT_SQL"))
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	return nil
}

func isAllowedTable(name string) bool {
	return AllowedTables[name]
}

func Create[T Helper](w http.ResponseWriter, r *http.Request) {
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}

	cols := item.GetNameColumns()         // e.g. "address, type, price, status"
	placeholders := item.GetPlaceholder() // e.g. "$1, $2, $3, $4"
	values := item.GetValues()

	// вычисляем следующий placeholder для owner_id
	// считаем количество placeholders (количество запятых + 1)
	n := 1
	if strings.TrimSpace(placeholders) != "" {
		n = len(strings.Split(placeholders, ","))
	}
	ownerPlaceholder := fmt.Sprintf("$%d", n+1)

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, owner_id) VALUES (%s, %s)",
		table,
		cols,
		placeholders,
		ownerPlaceholder,
	)

	ownerID, err := store.GetIDUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	args := append(values, ownerID)

	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Println("Create Exec error:", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Read[T Helper](w http.ResponseWriter, r *http.Request) {
	var item T
	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("SELECT * FROM %s", table)
	result := []T{}
	if err := DB.Select(&result, query); err != nil {
		log.Println("Read Select error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func Update[T Helper](w http.ResponseWriter, r *http.Request) {
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}

	param := strings.Split(item.GetNameColumns(), ", ")
	placeholder := strings.Split(item.GetPlaceholder(), ", ")
	if len(param) != len(placeholder) {
		http.Error(w, "Mismatch columns/placeholders", http.StatusInternalServerError)
		return
	}

	setParam := ""
	for i := 0; i < len(param); i++ {
		setParam += param[i] + " = " + strings.TrimSpace(placeholder[i]) + ", "
	}
	setParam = strings.TrimSuffix(setParam, ", ")

	// id placeholder должен быть следующим по номеру
	idPlaceholder := fmt.Sprintf("$%d", len(placeholder)+1)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %s",
		table,
		setParam,
		idPlaceholder,
	)

	args := append(item.GetValues(), id)
	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Println("Update Exec error:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Delete[T Helper](w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	var item T
	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1",
		table,
	)
	result, err := DB.Exec(query, id)
	if err != nil {
		log.Println("Delete Exec error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Операция завершилась успешно!"})
}

func GetByID[T Helper](w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	var item T
	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table)
	var result T
	if err := DB.Get(&result, query, id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Не найден!", http.StatusNotFound)
			return
		}
		log.Println("GetByID DB error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetMyData[T Helper](w http.ResponseWriter, r *http.Request) {
	userID, err := store.GetIDUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var item T
	table := item.GetNameTable()
	if !isAllowedTable(table) {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE owner_id = $1", table)
	result := []T{}
	if err := DB.Select(&result, query, userID); err != nil {
		log.Println("GetMyData DB error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
