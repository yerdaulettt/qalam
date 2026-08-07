package httphandler

import (
	"encoding/json"
	"net/http"

	"user-service/internal/auth"
)

type authHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(s *auth.AuthService) *authHandler {
	return &authHandler{service: s}
}

func (h *authHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userReg auth.RegisterReq

	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		errorResponse(w, err)
		return
	}

	user, err := h.service.Register(r.Context(), userReg)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ul auth.UserLogin

	err := json.NewDecoder(r.Body).Decode(&ul)
	if err != nil {
		errorResponse(w, err)
		return
	}

	tokens, err := h.service.Login(r.Context(), ul)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&tokens)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *authHandler) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var refresh struct {
		Refresh string `json:"refresh"`
	}

	err := json.NewDecoder(r.Body).Decode(&refresh)
	if err != nil {
		errorResponse(w, err)
		return
	}

	access, err := h.service.TokenRefresh(r.Context(), refresh.Refresh)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.Write([]byte(`{"access":"` + access + `"}`))
}
