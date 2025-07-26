package main

import (
	"subscriptions/rest-service/internal/api/handlers"
	"subscriptions/rest-service/internal/api/routers"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/service"
	"subscriptions/rest-service/pkg/database"
	"log"
	_ "subscriptions/rest-service/docs"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("config/")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("BD_HOST", "localhost")
	viper.SetDefault("BD_PORT", "5432")
	viper.SetDefault("BD_NAME", "test")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("\033[31merror reading .env file: %v\033[0m", err)
	}

	pass, user := viper.GetString("DB_PASSWORD"), viper.GetString("DB_USER")
	if pass == "" || user == "" {
		log.Fatalf("\033[31myou forgot set password or user for database!\033[0m")
	}
}

// @title           Subscription API With Swagger
// @version         1.0

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	db, err := database.GetDBConnect()
	if err != nil {
		log.Fatalf("\033[31merror connect to db: %v\033[0m", err)
	}

	subsRepo := repository.NewRepository(db)
	subsService := service.NewService(subsRepo)
	subsHandler := handlers.NewHandler(subsService)

	router := routers.SetupRouter(subsHandler)

	if err = router.Run("localhost:8080"); err != nil {
		log.Panicf("error start server: %v", err)
	}
}
