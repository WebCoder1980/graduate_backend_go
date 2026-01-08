package handler

import (
	"log"
	"net/http"
	"os"
)

const prefix = "/api/v1/user"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("User microservice handler test"))
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) Start() {
	http.HandleFunc(prefix, h.UserHandler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}
