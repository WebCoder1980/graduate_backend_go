package app

import (
	"graduate_backend_image_microservice_go/internal/handler"
	"graduate_backend_image_microservice_go/internal/kafka"
	"sync"
)

func Run() {
	var wg sync.WaitGroup

	wg.Go(kafka.Consumer)
	wg.Go(handler.Handler)

	wg.Wait()
}
