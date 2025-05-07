package main

import (
	"log"
	"net/http"
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

	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// func setupRoutes(H handlers.Handler) *chi.Mux {
// 	r := chi.NewRouter()

// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	r.Route("/people", func(r chi.Router) {
// 		r.Post("/", H.CreatePerson)
// 		r.Get("/", H.GetPeople)
// 		r.Put("/{id}", H.UpdatePerson)
// 		r.Delete("/{id}", H.DeletePerson)

// 	})

// 	r.Get("/swagger/*", httpSwagger.WrapHandler)

// 	return r
// }
