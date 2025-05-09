
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
	DeletePerson(ctx context.Context, id int64) error
}

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

func (uc *PersonUseCase) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCase) DeletePerson(ctx context.Context, id uint) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}


[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "DuplicateDecl",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "DuplicateDecl"
		}
	},
	"severity": 8,
	"message": "PersonUseCase redeclared in this block (see details)",
	"source": "compiler",
	"startLineNumber": 8,
	"startColumn": 6,
	"endLineNumber": 8,
	"endColumn": 19,
	"relatedInformation": [
		{
			"startLineNumber": 15,
			"startColumn": 6,
			"endLineNumber": 15,
			"endColumn": 19,
			"message": "",
			"resource": "/home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go"
		}
	]
}]

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "DuplicateDecl",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "DuplicateDecl"
		}
	},
	"severity": 8,
	"message": "PersonUseCase redeclared in this block",
	"source": "compiler",
	"startLineNumber": 15,
	"startColumn": 6,
	"endLineNumber": 15,
	"endColumn": 19,
	"relatedInformation": [
		{
			"startLineNumber": 8,
			"startColumn": 6,
			"endLineNumber": 8,
			"endColumn": 19,
			"message": "other declaration of PersonUseCase",
			"resource": "/home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go"
		}
	]
}]


