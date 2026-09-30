package main

import (
	"log"

	"system_components/internal/api"
	"system_components/internal/app/handler"
	"system_components/internal/app/repository"
)

func main() {
	repo := repository.NewRepository()
	h := handler.NewHandler(repo)

	router := api.SetupRouter(h)

	log.Println("Сервер слушает на http://localhost:8080/system-components")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}