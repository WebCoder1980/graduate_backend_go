package service

import (
	"graduate_backend_image_microservice_go/internal/kafka"
	"mime/multipart"
)

func Post(file multipart.File, filename string) {
	kafka.Producer(file, filename)
}
