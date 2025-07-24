package main

import (
	"subscriptions/rest-service/internal/api/handlers"
	"subscriptions/rest-service/internal/api/routers"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/service"
	"subscriptions/rest-service/pkg/database"
	"log"
	_ "subscriptions/rest-service/docs"
)

// @title           Subscription API With Swagger
// @version         1.0

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	db, err := database.GetDBConnect()
	if err != nil {
		log.Fatalf("error connect to db: %v", err)
	}

	subsRepo := repository.NewRepository(db)
	subsService := service.NewService(subsRepo)
	subsHandler := handlers.NewHandler(subsService)

	router := routers.SetupRouter(subsHandler)

	if err = router.Run("localhost:8080"); err != nil {
		log.Panicf("error start server: %v", err)
	}
}
