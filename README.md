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
	"context"
	"net/http"
	"os"

	"people/docs"
	"people/internal/db"
	"people/internal/handlers"
	"people/internal/logger"
	"people/internal/router"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Инициализация логгера с уровнем Debug (включает Debug, Info, Warn, Error)
	logger.SetLogger(logger.New(zapcore.DebugLevel))

	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found")
	} else {
		logger.Debug("Successfully loaded .env file")
	}

	// Получаем порт из переменных окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // fallback
	}

	// Логируем полученный порт на уровне Debug
	logger.Debugf("Using port: %s", port)

	// Инициализация базы данных
	database := db.Init()
	h := handlers.Handler{DB: database}

	// Инициализация маршрутов
	r := router.SetupRoutes(h)

	// Логируем запуск сервера
	ctx := logger.ToContext(context.Background(), logger.Global())
	logger.Infof(ctx, "Starting server on port: %s", port)

	// Логируем успешный запуск
	logger.Debug(ctx, "Routes setup completed")

	// Запуск сервера
	if err := http.ListenAndServe(":"+port, r); err != nil {
		// Логируем ошибку при запуске сервера
		logger.Fatalf(ctx, "Server failed: %v", err)
	} else {
		// Логируем успешный запуск сервера
		logger.Debug(ctx, "Server started successfully")
	}
}

[{
	"resource": "/home/gaz358/myprog/sobes/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use \"No .env file found\" (constant of type string) as context.Context value in argument to logger.Error: string does not implement context.Context (missing method Deadline)",
	"source": "compiler",
	"startLineNumber": 28,
	"startColumn": 16,
	"endLineNumber": 28,
	"endColumn": 36
}]






