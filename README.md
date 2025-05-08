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



package adapterhttp

import (
	"people/internal/app/usecase"
)

type Handler struct {
	CreateUC *usecase.CreatePersonUseCase
	// Здесь позже можно добавить другие usecase, если нужно
}

func NewHandler(createUC *usecase.CreatePersonUseCase) Handler {
	return Handler{CreateUC: createUC}
}
package adapterhttp

import (
	"encoding/json"
	"net/http"
	"people/internal/app/people"
)

func (h Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	person := people.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	createdPerson, err := h.CreateUC.Execute(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func (h Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}

func (h Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UpdatePerson not implemented yet"))
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeletePerson not implemented yet"))
}


package main

import (
	"log"
	"net/http"
	"people/internal/app/repository"
	"people/internal/app/services"
	"people/internal/app/usecase"
	"people/internal/adapterhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=qwert dbname=people port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	db.AutoMigrate(&people.Person{})

	repo := repository.NewPersonRepository(db)
	extService := services.NewExternalService() // реализуй этот сервис
	createUC := usecase.NewCreatePersonUseCase(repo, extService)
	handler := adapterhttp.NewHandler(createUC)

	r := adapterhttp.SetupRoutes(handler)
	log.Println("server started on :8080")
	http.ListenAndServe(":8080", r)
}



