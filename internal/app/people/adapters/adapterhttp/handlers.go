package adapterhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/repos"
	"people/internal/app/people/usecase"
	"people/pkg/logger"
	"strconv"

	"github.com/go-chi/chi"
)

type Handler struct {
	CreateUC   *usecase.CreatePersonUseCase
	PersonRepo repos.PersonRepository // Добавлено поле для репозитория
}

func NewHandler(createUC *usecase.CreatePersonUseCase) Handler {
	return Handler{CreateUC: createUC}
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

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	logger.Debug(ctx, "Delete request received", "id", idStr)

	if idStr == "" {
		logger.Warn(ctx, "No ID provided in URL")
		http.Error(w, `{"error":"missing ID"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn(ctx, "Invalid ID format", "id", idStr, "err", err)
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("DeletePerson ID=%d\n", id)

	_, err = h.PersonRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
		return
	}

	// if err := h.PersonRepo.Delete(ctx, id); err != nil {
	// 	logger.Error(ctx, "Failed to delete person", "id", id, "err", err)
	// 	http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
	// 	return
	// }

	logger.Info(ctx, "Person deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPeople not implemented yet"))
}

func (h Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UpdatePerson not implemented yet"))
}

// func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("DeletePerson not implemented yet"))
// }
