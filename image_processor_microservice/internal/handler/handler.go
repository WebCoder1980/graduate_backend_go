package handler

import (
	"context"
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

func (h *Handler) ImageIdGet(w http.ResponseWriter, r *http.Request) {
	token, err := h.getToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	imageId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Panic(err)
	}

	data, err := h.service.ImageGetById(imageId, &token)
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write(data)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) ImageIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ImageIdGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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
