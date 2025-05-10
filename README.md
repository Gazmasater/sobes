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


// GetPeople godoc
// @Summary      Получить список людей
// @Description  Возвращает список людей с фильтрацией, сортировкой и пагинацией
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        gender       query     string  false  "Пол"
// @Param        nationality  query     string  false  "Национальность"
// @Param        name         query     string  false  "Имя (поиск по подстроке)"
// @Param        surname      query     string  false  "Фамилия (поиск по подстроке)"
// @Param        patronymic   query     string  false  "Отчество (поиск по подстроке)"
// @Param        age          query     int     false  "Возраст (точное совпадение)"
// @Param        sort_by      query     string  false  "Поле сортировки (id, name, surname, patronymic, age, gender, nationality)"
// @Param        order        query     string  false  "Направление сортировки (asc, desc)"
// @Param        limit        query     int     false  "Количество записей (по умолчанию 10)"
// @Param        offset       query     int     false  "Смещение (для пагинации)"
// @Success      200  {array}   people.Person
// @Failure      500  {object}  map[string]string
// @Router       /people [get]





