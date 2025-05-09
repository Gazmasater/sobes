
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

repo := repository.NewPostgresRepository(db) // твоя реализация
createUC := usecase.NewCreatePersonUseCase(repo)
deleteUC := usecase.NewDeletePersonUseCase(repo)
personUC := usecase.NewPersonUseCase(createUC, deleteUC)


package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type CreatePersonUseCase struct {
	Repo repos.PersonRepository
}

func NewCreatePersonUseCase(repo repos.PersonRepository) *CreatePersonUseCase {
	return &CreatePersonUseCase{Repo: repo}
}

func (uc *CreatePersonUseCase) Execute(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.Repo.Create(ctx, person)
}

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id int64) error {
	return uc.Repo.Delete(ctx, id)
}

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error
}

type PersonUseCaseImpl struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
}

func NewPersonUseCase(
	createUC *CreatePersonUseCase,
	deleteUC *DeletePersonUseCase,
) *PersonUseCaseImpl {
	return &PersonUseCaseImpl{
		CreatePersonUseCase: createUC,
		DeletePersonUseCase: deleteUC,
	}
}

func (uc *PersonUseCaseImpl) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCaseImpl) DeletePerson(ctx context.Context, id int64) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}


package adapterhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"strconv"
)

type Handler struct {
	CreateUC usecase.PersonUseCase
	DeleteUC usecase.PersonUseCase // Добавляем новый UseCase для удаления
}

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

	fmt.Printf("PERSON NAme=%s Surname=%s\n", person.Name, person.Surname)

	createdPerson, err := h.CreateUC.Execute(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL

	idStr := r.URL.Path[len("/people/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.DeleteUC.Execute(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/adapters/adapterhttp/handlers.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "h.DeleteUC.Execute undefined (type usecase.PersonUseCase has no field or method Execute)",
	"source": "compiler",
	"startLineNumber": 57,
	"startColumn": 19,
	"endLineNumber": 57,
	"endColumn": 26
}]
