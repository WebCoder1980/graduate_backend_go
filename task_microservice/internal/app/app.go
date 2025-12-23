package app

import (
	"context"
	"graduate_backend_task_microservice/internal/handler"
	"log"
	"sync"
)

func Run() {
	ctx := context.Background()
	var wg sync.WaitGroup

	appHandler, err := handler.NewHandler(ctx)
	if err != nil {
		log.Panic(err)
	}
	wg.Go(appHandler.Start)

	wg.Wait()
}
