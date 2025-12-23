package app

import (
	"context"
	"graduate_backend_image_processor_microservice/internal/kafka"
	"log"
)

func Run() {
	ctx := context.Background()

	consumer, err := kafka.NewConsumer(ctx)
	if err != nil {
		log.Panic(err)
	}

	consumer.Start()
}
