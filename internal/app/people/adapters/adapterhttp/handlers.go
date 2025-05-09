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
	err = h.uc.DeletePerson(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Получаем тело запроса
	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем текущего человека по ID
	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	// Проверяем, изменилось ли имя
	nameChanged := existing.Name != req.Name

	// Обновляем поля
	existing.Name = req.Name
	existing.Surname = req.Surname
	existing.Patronymic = req.Patronymic

	if nameChanged {
		existing.Age = h.svc.GetAge(ctx, req.Name)
		existing.Gender = h.svc.GetGender(ctx, req.Name)
		existing.Nationality = h.svc.GetNationality(ctx, req.Name)
	}

	// Обновляем в базе через usecase
	updatedPerson, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(updatedPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {
	r.Post("/people", h.CreatePerson)
	r.Delete("/people/{id}", h.DeletePerson)
	r.Put("/people/{id}", h.UpdatePerson)

}

// func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("DeletePerson not implemented yet"))

// }

func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}

// func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("UpdatePerson not implemented yet"))
// }

// func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("DeletePerson not implemented yet"))
// }
