package user

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jadersonmarc/ecommerce-api/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Email    string
		Password string
	}

	json.NewDecoder(r.Body).Decode(&input)

	user, err := h.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string
		Password string
	}

	json.NewDecoder(r.Body).Decode(&input)

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, string(user.Role))

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	clamis := r.Context().Value(auth.UserContextKey).(jwt.MapClaims)

	userID := clamis["user_id"].(string)
	role := clamis["role"].(string)

	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"role":    role,
	})
}
