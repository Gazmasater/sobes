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


  curl -X GET http://localhost:8080/people

  gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ go run .
{"lvl":"debug","ts":"2025-05-10T22:02:20.961+0300","msg":"Successfully loaded .env file"}
{"lvl":"debug","ts":"2025-05-10T22:02:20.965+0300","msg":"Using port: 8080"}
{"lvl":"info","ts":"2025-05-10T22:02:21.008+0300","msg":"Starting server on port: 8080"}
2025/05/10 22:02:42 "GET http://localhost:8080/people/3 HTTP/1.1" from [::1]:59528 - 405 0B in 47.854µs








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




func AddIndexes(db *gorm.DB) error {
	// Добавление индекса для поля Name
	if err := db.Model(&Person{}).AddIndex("idx_name", "name").Error; err != nil {
		return err
	}

	// Добавление индекса для поля Surname
	if err := db.Model(&Person{}).AddIndex("idx_surname", "surname").Error; err != nil {
		return err
	}

	// Добавление индекса для поля Patronymic
	if err := db.Model(&Person{}).AddIndex("idx_patronymic", "patronymic").Error; err != nil {
		return err
	}

	// При необходимости добавьте индексы для других полей

	return nil
}

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/domain.go",
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
	"message": "db.Model(&Person{}).AddIndex undefined (type *gorm.DB has no field or method AddIndex)",
	"source": "compiler",
	"startLineNumber": 29,
	"startColumn": 32,
	"endLineNumber": 29,
	"endColumn": 40
}]
