
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


func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&people.Person{}, id).Error
}


type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}


type HTTPHandler struct {
	uc usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	// парсинг запроса и вызов h.uc.Create(...)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// парсинг запроса и вызов h.uc.Delete(...)
}


handler := adapterhttp.NewHandler(personUC) // <- теперь корректно



