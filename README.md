
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

  

package adapterhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"people/internal/serv"
	"strconv"

	"github.com/go-chi/chi"
)

type HTTPHandler_interf interface {
	RegisterRoutes(r chi.Router)
}

type HTTPHandler struct {
	svc serv.ExternalService

	uc usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase, svc serv.ExternalService) HTTPHandler_interf {
	return &HTTPHandler{uc: uc, svc: svc}
}

// CreatePerson godoc
// @Summary      Create person
// @Description  Creates a new person with enriched data
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        person  body      CreatePersonRequest  true  "Person to create"
// @Success      200     {object}  PersonResponse
// @Failure      400     {string}  string  "invalid request body"
// @Failure      500     {string}  string  "internal server error"
// @Router       /people [post]
func (h HTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем доп. данные
	ctx := r.Context()
	age := h.svc.GetAge(ctx, req.Name)
	gender := h.svc.GetGender(ctx, req.Name)
	nationality := h.svc.GetNationality(ctx, req.Name)

	person := people.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdPerson, err := h.uc.CreatePerson(ctx, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeletePerson godoc
// @Summary      Delete person
// @Description  Deletes person by ID
// @Tags         people
// @Produce      json
// @Param        id   path      int64  true  "Person ID"
// @Success      204  {string}  string  "no content"
// @Failure      400  {string}  string  "invalid id"
// @Failure      500  {string}  string  "internal server error"
// @Router       /people/{id} [delete]
func (h HTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL

	idStr := r.URL.Path[len("/people/"):]

	fmt.Printf("DeletePerson URL=%s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.uc.DeletePerson(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePerson godoc
// @Summary      Update person
// @Description  Updates person by ID and enriches if name changed
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        id      path      int64             true  "Person ID"
// @Param        person  body      PersonResponse    true  "Updated person"
// @Success      200     {object}  PersonResponse
// @Failure      400     {string}  string  "invalid request body or id"
// @Failure      404     {string}  string  "person not found"
// @Failure      500     {string}  string  "failed to update person"
// @Router       /people/{id} [put]
func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req PersonResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	nameChanged := existing.Name != req.Name

	existing.Name = req.Name
	existing.Surname = req.Surname
	existing.Patronymic = req.Patronymic
	existing.Age = req.Age
	existing.Gender = req.Gender
	existing.Nationality = req.Nationality

	if nameChanged {
		existing.Age = h.svc.GetAge(ctx, req.Name)
		existing.Gender = h.svc.GetGender(ctx, req.Name)
		existing.Nationality = h.svc.GetNationality(ctx, req.Name)
	}

	updated, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		http.Error(w, "failed to update person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(updated))
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {

	r.Post("/people", h.CreatePerson)
	r.Delete("/people/{id}", h.DeletePerson)
	r.Put("/people/{id}", h.UpdatePerson)

}

func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}




