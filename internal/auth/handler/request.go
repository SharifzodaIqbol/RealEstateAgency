package handler

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
