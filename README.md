
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


[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go",
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
	"message": "cannot use id (variable of type int64) as uint value in argument to uc.Repo.Delete",
	"source": "compiler",
	"startLineNumber": 30,
	"startColumn": 29,
	"endLineNumber": 30,
	"endColumn": 31
}]

