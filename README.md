people/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   ├── handlers/
│   └── router/
│       └── router.go


go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.

package main

import (
	"log"
	"net/http"

	"people/docs" // путь к swagger docs
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

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


4. ✏️ Пример аннотации в обработчике
📁 internal/handlers/person.go

// CreatePerson godoc
// @Summary      Create a new person
// @Description  Add person by JSON
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        person  body  models.Person  true  "Person"
// @Success      201     {object}  models.Person
// @Failure      400     {object}  map[string]string
// @Router       /people [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	// ...
}
5. 🔁 Обнови Swagger при изменении кода
Каждый раз после изменения аннотаций:


swag init
✅ Swagger будет доступен по адресу:

http://localhost:8080/swagger/index.html








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


gaz358@gaz358-BOD-WXX9:~/myprog/test/cmd$ swag init
2025/05/05 15:39:36 Generate swagger docs....
2025/05/05 15:39:36 Generate general API Info, search dir:./
2025/05/05 15:39:36 create docs.go at docs/docs.go
2025/05/05 15:39:36 create swagger.json at docs/swagger.json
2025/05/05 15:39:36 create swagger.yaml at docs/swagger.yaml
gaz358@gaz358-BOD-WXX9:~/myprog/test/cmd$ 





