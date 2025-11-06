package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscriptions/rest-service/docs"
	"subscriptions/rest-service/internal/api/handlers"
	"subscriptions/rest-service/internal/api/routers"
	"subscriptions/rest-service/internal/repository"
	"subscriptions/rest-service/internal/service"
	"subscriptions/rest-service/pkg/database"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("DB_HOST", "main_db")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("POSTGRES_DB", "test")
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("APP_PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("\033[31merror reading .env file: %v\033[0m", err)
	}

	pass, user := viper.GetString("POSTGRES_PASSWORD"), viper.GetString("POSTGRES_USER")
	if pass == "" || user == "" {
		log.Fatalf("\033[31myou forgot set password or user for database!\033[0m")
	}
}

// @title           Subscription API With Swagger
// @version         1.0

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	address := viper.GetString("APP_HOST") + ":" + viper.GetString("APP_PORT")

	docs.SwaggerInfo.Host = "localhost:" + viper.GetString("APP_PORT")

	db, err := database.GetDBConnect()
	if err != nil {
		log.Fatalf("\033[31merror connect to db: %v\033[0m", err)
	}

	subsRepo := repository.NewRepository(db)
	subsService := service.NewService(subsRepo)
	subsHandler := handlers.NewHandler(subsService)

	router := routers.SetupRouter(subsHandler)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server at %s\n", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("\033[31merror starting server: %v\033[0m", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("\033[31mserver forced to shutdown: %v\033[0m", err)
	}

	log.Println("Server stopped gracefully")
}
