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
	"os"
	"strconv"
	"time"

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
		port = "8080"
	}

	// ReadTimeout
	readTimeoutSec := 10 // default
	if readTimeoutStr := os.Getenv("READ_TIMEOUT"); readTimeoutStr != "" {
		if v, err := strconv.Atoi(readTimeoutStr); err == nil {
			readTimeoutSec = v
		} else {
			log.Printf("Invalid READ_TIMEOUT: %v", err)
		}
	}

	// WriteTimeout
	writeTimeoutSec := 10 // default
	if writeTimeoutStr := os.Getenv("WRITE_TIMEOUT"); writeTimeoutStr != "" {
		if v, err := strconv.Atoi(writeTimeoutStr); err == nil {
			writeTimeoutSec = v
		} else {
			log.Printf("Invalid WRITE_TIMEOUT: %v", err)
		}
	}

	database := db.Init()
	h := handlers.Handler{DB: database}
	r := router.SetupRoutes(h)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  time.Duration(readTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSec) * time.Second,
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(srv.ListenAndServe())
}






