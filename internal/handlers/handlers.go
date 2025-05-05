package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"people/internal/models"
	"people/internal/services"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var p models.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Gender = services.GetGender(p.Name)
	p.Age = services.GetAge(p.Name)
	p.Nationality = services.GetNationality(p.Name)

	h.DB.Create(&p)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person
	query := h.DB

	gender := r.URL.Query().Get("gender")
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	nationality := r.URL.Query().Get("nationality")
	if nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}

	query.Limit(limit).Offset(offset).Find(&people)
	json.NewEncoder(w).Encode(people)
}

func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Person
	if err := h.DB.First(&p, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	var updated models.Person
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.DB.Model(&p).Updates(updated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	h.DB.Delete(&models.Person{}, id)
	w.WriteHeader(http.StatusNoContent)
}
