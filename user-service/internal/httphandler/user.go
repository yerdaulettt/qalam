package httphandler

import (
	"encoding/json"
	"net/http"

	"user-service/internal/user"
)

type userHandler struct {
	service *user.UserService
}

func newUserHandler(s *user.UserService) *userHandler {
	return &userHandler{service: s}
}

func (h *userHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userReg user.RegisterReq

	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		errorResponse(w, err)
		return
	}

	u, err := h.service.Register(r.Context(), userReg)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&u)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ul user.UserLogin

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

func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		errorResponse(w, errUserId)
		return
	}

	u, err := h.service.GetMyProfile(r.Context(), userId)
	if err != nil {
		errorResponse(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&u)
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func (h *userHandler) TokenRefresh(w http.ResponseWriter, r *http.Request) {
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
