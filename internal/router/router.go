package router

import (
	"people/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(h handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/people", func(r chi.Router) {
		r.Post("/", h.CreatePerson)
		r.Get("/", h.GetPeople)
		r.Put("/{id}", h.UpdatePerson)
		r.Delete("/{id}", h.DeletePerson)
	})

	return r
}
