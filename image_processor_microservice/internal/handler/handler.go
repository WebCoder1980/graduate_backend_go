package handler

import (
	"context"
	"encoding/json"
	"errors"
	"graduate_backend_image_processor_microservice/internal/service"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const prefix = "/api/v1/image-processor"

type Handler struct {
	service *service.Service
}

func NewHandler(ctx context.Context) (*Handler, error) {
	serv, err := service.NewService(ctx)
	if err != nil {
		return nil, err
	}

	return &Handler{service: serv}, nil
}

func NewHandlerWithService(serv *service.Service) *Handler {
	return &Handler{service: serv}
}

func (h *Handler) ImageIdGet(w http.ResponseWriter, r *http.Request) {
	token, err := h.getToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	imageId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || imageId <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	data, err := h.service.ImageGetById(imageId, &token)
	if err != nil {
		log.Printf("ImageIdGet error: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to get image")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) ImageIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.ImageIdGet(w, r)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handler) Start() {
	http.HandleFunc(prefix+"/{id}", h.ImageIdHandler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}

func (h *Handler) getToken(r *http.Request) (string, error) {
	val := r.Header.Values("authorization")

	if len(val) == 0 {
		return "", errors.New("The 'authorization' header is missing")
	}

	res, found := strings.CutPrefix(val[0], "Bearer ")
	if !found {
		return "", errors.New("The 'authorization' header don't have prefix 'Bearer '")
	}

	return res, nil
}
