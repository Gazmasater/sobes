
DROP TABLE IF EXISTS people;


internal/
└── app/
    └── mydomain/
        ├── usecase/
        │   ├── user_usecase.go        # Бизнес-логика
        │   └── user_usecase_iface.go  # Интерфейс, например UserRepository
        ├── repository/
        │   └── postgres/
        │       └── user_repository.go# Реализация интерфейса
        ├── adapters/
        │   └── http/
        │       └── handler.go         # Использует интерфейс Usecase
        └── domain.go


 curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ivan",
    "surname": "Seli",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/people/1"


curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович",
    "age": 35,
    "gender": "male",
    "nationality": "russian"
  }'






go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.


git rm --cached textDB


curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dmitriy",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'

  curl -X DELETE "http://localhost:8080/people/26"


  curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alexey",
    "surname": "Ivanov",
    "patronymic": "Sergeevich",
    "age": 30,
    "gender": "male",
    "nationality": "ru"
  }'

  



package main

import (
	"context"
	"log"
	"net/http"
	"os"
	_ "people/docs"
	"people/pkg/logger"

	"people/internal/app/people"
	"people/internal/app/people/adapters/adapterhttp"
	"people/internal/app/people/repos"
	"people/internal/app/people/usecase"
	"people/internal/serv"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	logger.SetLogger(logger.New(zapcore.DebugLevel))

	ctx := logger.ToContext(context.Background(), logger.Global())
	r := chi.NewRouter()

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port)

	dsn := "host=localhost user=postgres password=qwert dbname=people port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	// Миграция таблицы Person
	db.AutoMigrate(&people.Person{})

	// Создание зависимостей
	repo := repos.NewPersonRepository(db)

	// Create and Delete UseCases
	createUC := usecase.NewCreatePersonUseCase(repo)
	deleteUC := usecase.NewDeletePersonUseCase(repo)

	// Объединённый интерфейс
	personUC := usecase.NewPersonUseCase(createUC, deleteUC)
	svc := serv.NewExternalService()

	// Handler принимает один интерфейс
	handler := adapterhttp.NewHandler(personUC, svc)

	handler.RegisterRoutes(r)

	// Запуск сервера
	log.Printf("server started on :%s", port)
	http.ListenAndServe(":8080", r)
}



