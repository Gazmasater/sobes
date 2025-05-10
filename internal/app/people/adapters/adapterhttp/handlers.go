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
	"github.com/go-chi/chi/middleware"
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
// @Param        id      path      int64                 true  "Person ID"
// @Param        person  body      UpdatePersonRequest   true  "Updated person (partial)"
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

	// Декодируем JSON-тело запроса
	var req UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем существующего человека
	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	// Проверяем и обновляем только переданные поля
	if req.Name != nil {
		nameChanged := existing.Name != *req.Name
		existing.Name = *req.Name
		if nameChanged {
			existing.Age = h.svc.GetAge(ctx, *req.Name)
			existing.Gender = h.svc.GetGender(ctx, *req.Name)
			existing.Nationality = h.svc.GetNationality(ctx, *req.Name)
		}
	}
	if req.Surname != nil {
		existing.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		existing.Patronymic = *req.Patronymic
	}
	if req.Age != nil {
		existing.Age = *req.Age
	}
	if req.Gender != nil {
		existing.Gender = *req.Gender
	}
	if req.Nationality != nil {
		existing.Nationality = *req.Nationality
	}

	// Обновляем запись
	updated, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		http.Error(w, "failed to update person", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(updated))
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/people", h.CreatePerson)
	r.Delete("/people/{id}", h.DeletePerson)
	r.Put("/people/{id}", h.UpdatePerson)

}

func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}
