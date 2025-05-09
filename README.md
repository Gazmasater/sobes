
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
    "surname": "Selivanov",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/person/26"





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


package usecase

import (
	"context"
	"people/internal/app/people"
)

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id uint) error
}


package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
	"people/internal/serv"
)

type PersonUseCase struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
}

func NewPersonUseCase(
	createUC *CreatePersonUseCase,
	deleteUC *DeletePersonUseCase,
) *PersonUseCase {
	return &PersonUseCase{
		CreatePersonUseCase: createUC,
		DeletePersonUseCase: deleteUC,
	}
}

// Реализация методов интерфейса для создания и удаления
func (uc *PersonUseCase) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCase) DeletePerson(ctx context.Context, id uint) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}


package handlers

import (
	"context"
	"fmt"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
)

type Handler struct {
	PersonUseCase usecase.PersonUseCase
}

func NewHandler(personUseCase usecase.PersonUseCase) *Handler {
	return &Handler{
		PersonUseCase: personUseCase,
	}
}

// Обработчик для создания
func (h *Handler) CreatePersonHandler(ctx context.Context, person people.Person) (people.Person, error) {
	return h.PersonUseCase.CreatePerson(ctx, person)
}

// Обработчик для удаления
func (h *Handler) DeletePersonHandler(ctx context.Context, id uint) error {
	return h.PersonUseCase.DeletePerson(ctx, id)
}



package main

import (
	"context"
	"fmt"
	"log"
	"people/internal/app/people"
	"people/internal/app/people/adapters/adapterhttp/handlers"
	"people/internal/app/people/repos"
	"people/internal/app/people/usecase"
	"people/internal/serv"
	// другие импорты
)

func main() {
	// Инициализация внешнего сервиса (например, для получения данных о человеке)
	externalService := serv.NewExternalService()

	// Инициализация репозитория (например, подключение к базе данных)
	personRepo := repos.NewPersonRepository()

	// Инициализация UseCase для создания
	createPersonUseCase := usecase.NewCreatePersonUseCase(personRepo, externalService)

	// Инициализация UseCase для удаления
	deletePersonUseCase := usecase.NewDeletePersonUseCase(personRepo)

	// Инициализация общего UseCase для создания и удаления
	personUseCase := usecase.NewPersonUseCase(createPersonUseCase, deletePersonUseCase)

	// Инициализация обработчика HTTP
	handler := handlers.NewHandler(personUseCase)

	// Пример использования: создание персоны
	person := people.Person{
		Name:      "John",
		Surname:   "Doe",
		Patronymic: "Middle",
	}

	// Создание персоны
	createdPerson, err := handler.CreatePersonHandler(context.Background(), person)
	if err != nil {
		log.Fatalf("Error creating person: %v", err)
	}
	fmt.Printf("Created person: %+v\n", createdPerson)

	// Пример использования: удаление персоны по ID
	err = handler.DeletePersonHandler(context.Background(), createdPerson.ID)
	if err != nil {
		log.Fatalf("Error deleting person: %v", err)
	}
	fmt.Println("Person deleted successfully")
}



