package handler

import (
	"auth-service/internal/auth/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if req.Username == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Требуется имя пользователя"})
		return
	}
	if req.Email == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Требуется email пользователя"})
		return
	}
	if len(req.Password) < 6 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Пароль должен быть длиной не менее 6 символов."})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное тело запроса"})
		return
	}

	err := h.authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное тело запроса"})
		return
	}

	accessToken, refreshToken, err := h.authService.LoginUser(req.Identifier, req.Password)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}
