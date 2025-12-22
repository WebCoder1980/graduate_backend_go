package app

import (
	"context"
	"graduate_backend_task_microservice/internal/handler"
	"graduate_backend_task_microservice/internal/kafka"
	"log"
	"sync"
)

func Run() {
	ctx := context.Background()
	var wg sync.WaitGroup

	kafkaConsumer, err := kafka.NewConsumer(ctx)
	if err != nil {
		log.Panic(err)
	}
	wg.Go(kafkaConsumer.Start)

	appHandler, err := handler.NewHandler(ctx)
	if err != nil {
		log.Panic(err)
	}
	wg.Go(appHandler.Start)

	wg.Wait()
}
