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


# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=qwert
DB_NAME=people
DB_SSLMODE=disable

# Внешние API
AGIFY_API=https://api.agify.io
GENDERIZE_API=https://api.genderize.io
NATIONALIZE_API=https://api.nationalize.io


func Init() *gorm.DB {
	// Загрузка .env (если ещё не загружено)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Автоматическая миграция
	if err := db.AutoMigrate(&models.Person{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}





