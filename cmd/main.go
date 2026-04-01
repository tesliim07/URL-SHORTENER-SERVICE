package main

import (
	"net/http"
	"os"
	"fmt"
	
	"url-shortener-service/config"
	"url-shortener-service/internal/cache"
	"url-shortener-service/internal/repository"
	"url-shortener-service/internal/service"
	"url-shortener-service/internal/handler"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-chi/chi/v5"
 	"github.com/go-chi/chi/v5/middleware"
 	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// set up zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load(".env")
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	config := config.LoadConfig()
	//connect to database PostgreSQL
	repo, err := repository.NewRepo(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	//connect to Redis
	redisClient, err := cache.NewCache(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	service := service.NewService(repo, redisClient, config)
	handler := handler.NewHandler(service)

	router := chi.NewRouter()

	//middlewares
	router.Use(middleware.Logger) //Logs every request

	//Todo routes
	router.Get("/{code}", handler.Redirect)
	router.Post("/shorten", handler.ShortenURL)

	//Swagger route
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}