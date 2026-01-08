package app

import "user_microservice/internal/handler"

func Run() {
	hand := handler.NewHandler()
	hand.Start()
}
