package handler

import (
	"context"
	"encoding/json"
	"graduate_backend_task_microservice/internal/constant"
	"graduate_backend_task_microservice/internal/service"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.TaskPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(constant.FileMaxSize)
	if err != nil {
		log.Panic(err)
	}

	files := r.MultipartForm

	taskId, err := h.service.Post(files)
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(strconv.FormatInt(taskId, 10)))
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) TaskIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.TaskGetById(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskGetById(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panic(err)
	}

	result, err := h.service.GetImagesByTaskId(taskId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}

	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panic(err)
	}
}

func (h *Handler) Start() {
	http.HandleFunc(prefix, h.TaskHandler)
	http.HandleFunc(prefix+"/{id}", h.TaskIdHandler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}
