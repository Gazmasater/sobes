package main

import (
	"log"
	"net/http"
	"os"
	_ "people/docs"

	"people/internal/db"
	"people/internal/handlers"
	"people/internal/router"

	"github.com/joho/godotenv"
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // fallback
	}

	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	// Запуск сервера
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
