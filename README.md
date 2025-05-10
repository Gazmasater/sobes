go install github.com/swaggo/swag/cmd/swag@latest

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

  curl -X DELETE "http://localhost:8080/people/5"


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

  
swag init -g cmd/main.go -o docs




package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type CreatePersonUseCase struct {
	Repo repos.PersonRepository
}

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

type PersonUseCaseImpl struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
}

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error

	UpdatePerson(ctx context.Context, person people.Person) (people.Person, error)
	GetPersonByID(ctx context.Context, id int64) (people.Person, error)

	GetPeople(ctx context.Context) ([]people.Person, error)
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

func NewCreatePersonUseCase(repo repos.PersonRepository) *CreatePersonUseCase {
	return &CreatePersonUseCase{Repo: repo}
}

func (uc *CreatePersonUseCase) Execute(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.Repo.Create(ctx, person)
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id int64) error {
	return uc.Repo.Delete(ctx, id)
}

func (uc *PersonUseCaseImpl) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCaseImpl) DeletePerson(ctx context.Context, id int64) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}

func (uc *PersonUseCaseImpl) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.GetByID(ctx, id)
}

func (uc *PersonUseCaseImpl) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.Update(ctx, person)
}

