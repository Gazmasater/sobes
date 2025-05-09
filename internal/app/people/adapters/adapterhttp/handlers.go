package adapterhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"strconv"
)

type Handler struct {
	CreateUC *usecase.CreatePersonUseCase
	DeleteUC *usecase.DeletePersonUseCase // Добавляем новый UseCase для удаления
}

func NewHandler_C(createUC *usecase.CreatePersonUseCase) Handler {
	return Handler{CreateUC: createUC}
}

func NewHandler_D(deleteUC *usecase.DeletePersonUseCase) Handler {
	return Handler{DeleteUC: deleteUC}
}

func (h Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
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

	createdPerson, err := h.CreateUC.Execute(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL

	idStr := r.URL.Path[len("/people/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.DeleteUC.Execute(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("DeletePerson not implemented yet"))

// }

func (h Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}

func (h Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UpdatePerson not implemented yet"))
}

// func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("DeletePerson not implemented yet"))
// }
