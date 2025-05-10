
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

  
gaz358@gaz358-BOD-WXX9:~/myprog/sobes/cmd$ go run .
{"lvl":"debug","ts":"2025-05-10T18:30:43.147+0300","msg":"Successfully loaded .env file"}
{"lvl":"debug","ts":"2025-05-10T18:30:43.155+0300","msg":"Using port: 8081"}
{"lvl":"info","ts":"2025-05-10T18:30:43.252+0300","msg":"Starting server on port: 8081"}
2025/05/10 18:30:47 "GET http://localhost:8081/swagger/index.html HTTP/1.1" from 127.0.0.1:54304 - 200 3373B in 613.675µs
2025/05/10 18:30:47 "GET http://localhost:8081/swagger/swagger-ui.css HTTP/1.1" from 127.0.0.1:54304 - 200 143980B in 525.731µs
2025/05/10 18:30:47 "GET http://localhost:8081/swagger/swagger-ui-bundle.js HTTP/1.1" from 127.0.0.1:54310 - 200 1095116B in 890.541µs
2025/05/10 18:30:47 "GET http://localhost:8081/swagger/swagger-ui-standalone-preset.js HTTP/1.1" from 127.0.0.1:54304 - 200 339540B in 492.562µs
2025/05/10 18:30:47 "GET http://localhost:8081/swagger/doc.json HTTP/1.1" from 127.0.0.1:54310 - 500 22B in 23.143µs









