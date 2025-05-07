package router

import (
	"people/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(h handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/people", func(r chi.Router) {
		r.Post("/", h.CreatePerson)
		r.Get("/", h.GetPeople)
		r.Put("/{id}", h.UpdatePerson)
		r.Delete("/{id}", h.DeletePerson) // Используем прямую передачу параметра в обработчик
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
