
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

  

package adapterhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"people/internal/serv"
	"strconv"

	"github.com/go-chi/chi"
)

type HTTPHandler_interf interface {
	RegisterRoutes(r chi.Router)
}

type HTTPHandler struct {
	svc serv.ExternalService

	uc usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase, svc serv.ExternalService) HTTPHandler_interf {
	return &HTTPHandler{uc: uc, svc: svc}
}

func (h HTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем доп. данные
	ctx := r.Context()
	age := h.svc.GetAge(ctx, req.Name)
	gender := h.svc.GetGender(ctx, req.Name)
	nationality := h.svc.GetNationality(ctx, req.Name)

	person := people.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdPerson, err := h.uc.CreatePerson(ctx, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h HTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL

	idStr := r.URL.Path[len("/people/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.uc.DeletePerson(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req PersonResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	nameChanged := existing.Name != req.Name

	existing.Name = req.Name
	existing.Surname = req.Surname
	existing.Patronymic = req.Patronymic
	existing.Age = req.Age
	existing.Gender = req.Gender
	existing.Nationality = req.Nationality

	if nameChanged {
		existing.Age = h.svc.GetAge(ctx, req.Name)
		existing.Gender = h.svc.GetGender(ctx, req.Name)
		existing.Nationality = h.svc.GetNationality(ctx, req.Name)
	}

	updated, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		http.Error(w, "failed to update person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(updated))
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {

	r.Post("/people", h.CreatePerson)
	r.Delete("/people/{id}", h.DeletePerson)
	r.Put("/people/{id}", h.UpdatePerson)

}

func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}


package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "people/docs"
	"people/pkg/logger"
	"time"

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

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /
func main() {

	logger.SetLogger(logger.New(zapcore.DebugLevel))

	ctx := logger.ToContext(context.Background(), logger.Global())
	r := chi.NewRouter()

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
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
		log.Fatal("failed to connect to DB:", err)
	}

	port_s := os.Getenv("SERVER_PORT")
	if port_s == "" {
		port_s = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port_s)

	db.AutoMigrate(&people.Person{})

	repo := repos.NewPersonRepository(db)

	createUC := usecase.NewCreatePersonUseCase(repo)
	deleteUC := usecase.NewDeletePersonUseCase(repo)

	personUC := usecase.NewPersonUseCase(createUC, deleteUC)
	svc := serv.NewExternalService()

	handler := adapterhttp.NewHandler(personUC, svc)

	handler.RegisterRoutes(r)

	logger.Infof(ctx, "Starting server on port: %s", port_s)

	srv := &http.Server{
		Addr:         ":" + port_s,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf(ctx, "Server failed: %v", err)
	}
}





