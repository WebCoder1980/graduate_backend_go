package handler

import (
	"context"
	"encoding/json"
	"errors"
	"graduate_backend_task_microservice/internal/constant"
	"graduate_backend_task_microservice/internal/service"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const prefix = "/api/v1/task"

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

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.TaskGet(w, r)
	case http.MethodPost:
		h.TaskPost(w, r)
	}
}

func (h *Handler) TaskGet(w http.ResponseWriter, r *http.Request) {
	token, err := h.getToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	result, err := h.service.TaskGetByUserUuid(&token)
	if err != nil {
		log.Printf("TaskGet error: %v", err)
		if strings.Contains(err.Error(), "token") {
			writeError(w, http.StatusUnauthorized, "Invalid or missing token")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get tasks")
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) TaskPost(w http.ResponseWriter, r *http.Request) {
	if !validateContentType(r, "multipart/form-data") {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be multipart/form-data")
		return
	}

	token, err := h.getToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	err = r.ParseMultipartForm(constant.FileMaxSize)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	imageUrls := r.MultipartForm.Value["image_url"]
	if (r.MultipartForm == nil || len(r.MultipartForm.File) == 0) && len(imageUrls) == 0 {
		writeError(w, http.StatusBadRequest, "No files or URLs provided")
		return
	}

	var width, height *int
	var format *string
	var quality *float64

	query := r.URL.Query()

	if query.Has("width") {
		wVal, err := strconv.Atoi(query.Get("width"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid width parameter")
			return
		}
		width = &wVal
	}

	if query.Has("height") {
		hVal, err := strconv.Atoi(query.Get("height"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid height parameter")
			return
		}
		height = &hVal
	}

	if query.Has("format") {
		f := query.Get("format")
		format = &f
	}

	if query.Has("quality") {
		q, err := strconv.ParseFloat(query.Get("quality"), 32)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid quality parameter")
			return
		}
		quality = &q
	}

	taskId, err := h.service.Post(r.MultipartForm, imageUrls, width, height, format, quality, &token)
	if err != nil {
		log.Printf("TaskPost error: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(strconv.FormatInt(taskId, 10)))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) TaskIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.TaskGetById(w, r)
}

func (h *Handler) TaskGetById(w http.ResponseWriter, r *http.Request) {
	token, err := h.getToken(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	taskId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || taskId <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	result, err := h.service.GetImagesByTaskId(taskId, &token)
	if err != nil {
		log.Printf("TaskGetById error: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to get task images")
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (h *Handler) Start() {
	http.HandleFunc(prefix, h.TaskHandler)
	http.HandleFunc(prefix+"/{id}", h.TaskIdHandler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func validateContentType(r *http.Request, expected string) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, expected)
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
