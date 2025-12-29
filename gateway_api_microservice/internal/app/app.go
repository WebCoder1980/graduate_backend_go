package app

import "gateway_api_microservice/internal/handler"

func Run() {
	hand := handler.NewHandler()

	hand.Start()
}
