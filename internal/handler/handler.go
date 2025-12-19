package handler

import (
	"graduate_backend_image_microservice_go/internal/constant"
	"graduate_backend_image_microservice_go/internal/service"
	"log"
	"net/http"
)

const prefix = "/api/v1/image"

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(constant.FileMaxSize)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла: ", http.StatusBadRequest)
		return
	}
	defer file.Close()

	service.Post(file, handler.Filename)
}

func Handler() {
	http.HandleFunc(prefix+"/", Post)

	log.Fatal(http.ListenAndServe(":5267", nil))
}
