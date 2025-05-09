
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
  

package adapterhttp

import (
	"net/http"

	"people/internal/app/people/usecase"

	"github.com/go-chi/chi/v5"
)

type HTTPHandler interface {
	RegisterRoutes(r chi.Router)
}

type handler struct {
	uc usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase) HTTPHandler {
	return &handler{uc: uc}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Post("/person", h.CreatePerson)
	r.Delete("/person/{id}", h.DeletePerson)
}

func (h *handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Обработка создания
}

func (h *handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Обработка удаления
}



package adapterhttp

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetupRoutes(h HTTPHandler) http.Handler {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}



// ...
handler := adapterhttp.NewHandler(personUC) // <- возвращается интерфейс HTTPHandler
r := adapterhttp.SetupRoutes(handler)

log.Println("server started on :8080")
http.ListenAndServe(":8080", r)

