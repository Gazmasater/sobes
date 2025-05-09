
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
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"strconv"
)

type HTTPHandler struct {
	UseCase usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h HTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
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

	createdPerson, err := h.UseCase.CreatePerson(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h HTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL

	idStr := r.URL.Path[len("/people/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.UseCase.DeletePerson(r.Context(), int64(id))
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
		"value": "MissingLitField",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingLitField"
		}
	},
	"severity": 8,
	"message": "unknown field uc in struct literal of type HTTPHandler",
	"source": "compiler",
	"startLineNumber": 17,
	"startColumn": 22,
	"endLineNumber": 17,
	"endColumn": 24
}]

