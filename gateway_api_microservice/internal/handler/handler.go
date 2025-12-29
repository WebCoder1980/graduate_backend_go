package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) handler(w http.ResponseWriter, r *http.Request) {
	const prefix = "/api/v1"

	if strings.HasPrefix(r.URL.Path, prefix+"/task") {
		h.getProxy(os.Getenv("task_microservice_address")).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, prefix+"/image-processor") {
		h.getProxy(os.Getenv("image_processor_microservice_address")).ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *Handler) getProxy(target string) *httputil.ReverseProxy {
	res, err := url.Parse(target)
	if err != nil {
		log.Panic(err)
	}

	return httputil.NewSingleHostReverseProxy(res)
}

func (h *Handler) Start() {
	http.HandleFunc("/", h.handler)

	log.Panic(http.ListenAndServe(":"+os.Getenv("handler_port"), nil))
}
