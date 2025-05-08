package adapterhttp

import (
	"encoding/json"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
)

type Handler struct {
	CreateUC *usecase.CreatePersonUseCase
}

func NewHandler(createUC *usecase.CreatePersonUseCase) *Handler {
	return &Handler{CreateUC: createUC}
}

func (h *Handler) CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
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

	createdPerson, err := h.CreateUC.Execute(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ToResponse(createdPerson)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
