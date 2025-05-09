
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


package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id uint) error // Новый метод для удаления
}


package repos

import (
	"context"
	"people/internal/app/people"

	"gorm.io/gorm"
)

// GormPersonRepository реализация PersonRepository через GORM
type GormPersonRepository struct {
	db *gorm.DB
}

// NewPersonRepository создаёт новый GormPersonRepository
func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

// Create сохраняет нового человека в базу данных
func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

// Delete удаляет человека по ID
func (r *GormPersonRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}


package adapterhttp

import (
	"fmt"
	"net/http"
	"strconv"

	"people/internal/app/people/usecase"
)

type Handler struct {
	CreateUC   *usecase.CreatePersonUseCase
	DeleteUC   *usecase.DeletePersonUseCase // Добавляем новый UseCase для удаления
}

func NewHandler(createUC *usecase.CreatePersonUseCase, deleteUC *usecase.DeletePersonUseCase) Handler {
	return Handler{CreateUC: createUC, DeleteUC: deleteUC}
}

func (h Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Обработчик для создания человека
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	idStr := r.URL.Path[len("/persons/"):]
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


package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id uint) error {
	return uc.Repo.Delete(ctx, id)
}


deleteUC := usecase.NewDeletePersonUseCase(personRepo)
createUC := usecase.NewCreatePersonUseCase(personRepo)

handler := adapterhttp.NewHandler(createUC, deleteUC)



type Handler struct {
	CreateUC *usecase.CreatePersonUseCase
	DeleteUC *usecase.DeletePersonUseCase // Добавляем новый UseCase для удаления
}

func NewHandler_C(createUC *usecase.CreatePersonUseCase) Handler {
	return Handler{CreateUC: createUC}
}

func NewHandler_D(deleteUC *usecase.CreatePersonUseCase) Handler {
	return Handler{DeleteUC: deleteUC}
}

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/adapters/adapterhttp/handlers.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use deleteUC (variable of type *usecase.CreatePersonUseCase) as *usecase.DeletePersonUseCase value in struct literal",
	"source": "compiler",
	"startLineNumber": 22,
	"startColumn": 27,
	"endLineNumber": 22,
	"endColumn": 35
}]



