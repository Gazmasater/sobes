package main

import (
	"log"
	"net/http"
	"people/internal/db"
	"people/internal/handlers"
	"people/internal/router"
)

func main() {
	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
