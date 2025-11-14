package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	auth_model "auth-service/internal/auth/model"
	"auth-service/internal/estate/model"
	"auth-service/internal/estate/repository"

	"github.com/go-chi/chi/v5"
)

type contextKey string

const userContextKey contextKey = "agent"

type EstateHandler struct {
	repo repository.EstateRepository
}

func NewEstateHandler(repo repository.EstateRepository) *EstateHandler {
	return &EstateHandler{repo: repo}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type CreatePropertyRequest struct {
	Address string  `json:"address"`
	Type    string  `json:"type"`
	Price   float64 `json:"price"`
}

type CreateSaleRequest struct {
	PropertyID int     `json:"property_id"`
	BuyerID    int     `json:"buyer_id"`
	FinalPrice float64 `json:"final_price"`
}

type CreatePurchaseRequest struct {
	PropertyID   int     `json:"property_id"`
	SellerID     int     `json:"seller_id"`
	InitialPrice float64 `json:"initial_price"`
}

// CreateProperty создает новый объект недвижимости
func (h *EstateHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromContext(r)
	if err != nil {
		log.Printf("Auth error: %v", err)
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("Creating property for user ID: %d, Role: %s", claims.UserID, claims.RoleName)

	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	property := &model.Property{
		Address: req.Address,
		Type:    req.Type,
		Price:   req.Price,
		AgentID: claims.UserID,
		Status:  "available",
	}

	id, err := h.repo.CreateProperty(r.Context(), property)
	if err != nil {
		log.Printf("Error creating property: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create property"})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Property created successfully",
		"id":      fmt.Sprintf("%d", id),
	})
}

// GetPropertyByID возвращает объект недвижимости по ID
func (h *EstateHandler) GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid property ID"})
		return
	}

	property, err := h.repo.GetPropertyByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting property: %v", err)
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Property not found"})
		return
	}

	respondWithJSON(w, http.StatusOK, property)
}

func (h *EstateHandler) ListProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.repo.ListProperties(r.Context())
	if err != nil {
		log.Printf("Error listing properties: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve property list"})
		return
	}

	respondWithJSON(w, http.StatusOK, properties)
}

func (h *EstateHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth_model.CustomClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User claims not found"})
		return
	}

	var req CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	sale := &model.Sale{
		PropertyID: req.PropertyID,
		BuyerID:    req.BuyerID,
		FinalPrice: req.FinalPrice,
		AgentID:    claims.UserID,
	}

	id, err := h.repo.CreateSale(r.Context(), sale)
	if err != nil {
		log.Printf("Error creating sale: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to record sale"})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Sale recorded successfully",
		"id":      fmt.Sprintf("%d", id),
	})
}

func (h *EstateHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth_model.CustomClaims)
	if !ok {
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User claims not found"})
		return
	}

	var req CreatePurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	purchase := &model.Purchase{
		PropertyID:   req.PropertyID,
		SellerID:     req.SellerID,
		InitialPrice: req.InitialPrice,
		AgentID:      claims.UserID,
	}

	id, err := h.repo.CreatePurchase(r.Context(), purchase)
	if err != nil {
		log.Printf("Error creating purchase: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to record purchase"})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Purchase recorded successfully",
		"id":      fmt.Sprintf("%d", id),
	})
}

// Вспомогательная функция для получения claims из контекста
func getClaimsFromContext(r *http.Request) (*auth_model.CustomClaims, error) {
	user := r.Context().Value(userContextKey)
	if user == nil {
		return nil, fmt.Errorf("user claims not found in context")
	}

	claims, ok := user.(*auth_model.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid user claims type")
	}

	return claims, nil
}
