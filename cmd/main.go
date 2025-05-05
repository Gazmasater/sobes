package main

import (
	"log"
	"net/http"
	_ "people/docs"

	"people/internal/db"
	"people/internal/handlers"
	"people/internal/router"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /
func main() {
	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	r.Get("/swagger/", httpSwagger.WrapHandler)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
