package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"user_microservice/internal/model"
	"user_microservice/internal/service"
)

const prefix = "/api/v1/user"

type Handler struct {
	ctx     context.Context
	service *service.Service
}

func NewHandler(ctx context.Context) (*Handler, error) {
	serv, err := service.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &Handler{
		ctx:     ctx,
		service: serv,
	}, err
}

func (h *Handler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if !validateContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var body model.UserLogin
	if err = json.Unmarshal(data, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	res, err := h.service.UserLogin(&body)
	if err != nil {
		log.Printf("Login error: %v", err)
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if data, err = json.Marshal(res); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) UserRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if !validateContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var body model.UserRefreshToken
	if err := json.Unmarshal(data, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if body.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	result, err := h.service.UserRefreshToken(&body)
	if err != nil {
		log.Printf("Refresh token error: %v", err)
		writeError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if data, err = json.Marshal(result); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if !validateContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var body model.UserRegisterRequest
	if err = json.Unmarshal(data, &body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if body.Username == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	if len(body.Password) < 6 {
		writeError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	if len(body.Username) < 3 {
		writeError(w, http.StatusBadRequest, "Username must be at least 3 characters")
		return
	}

	err = h.service.UserRegisterPost(&body)
	if err != nil {
		log.Printf("Register error: %v", err)
		writeError(w, http.StatusConflict, "User already exists")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Start() {
	http.HandleFunc(prefix+"/login", h.UserLoginHandler)
	http.HandleFunc(prefix+"/refresh-token", h.UserRefreshTokenHandler)
	http.HandleFunc(prefix+"/register", h.UserRegisterHandler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func validateContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}
