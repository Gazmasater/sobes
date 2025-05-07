package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"people/internal/models"
	"people/internal/services"
)

type Handler struct {
	DB *gorm.DB
}

// CreatePerson godoc
// @Summary Создать нового человека
// @Description Принимает имя, фамилию и (опционально) отчество, автоматически определяет пол, возраст и национальность
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.CreatePersonRequest true "Данные для создания"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Router /people [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение данных из внешних API
	age := services.GetAge(req.Name)

	gender := services.GetGender(req.Name)

	nationality := services.GetNationality(req.Name)

	p := models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := h.DB.Create(&p).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetPeople godoc
// @Summary Получить список людей
// @Description Возвращает список людей с возможностью фильтрации по полу и национальности, а также с пагинацией
// @Tags people
// @Accept json
// @Produce json
// @Param gender query string false "Пол (например, male, female)"
// @Param nationality query string false "Национальность (например, Russian, American)"
// @Param limit query int false "Количество возвращаемых записей (по умолчанию 10)"
// @Param offset query int false "Смещение (offset) для пагинации"
// @Success 200 {array} models.Person
// @Failure 500 {object} map[string]string
// @Router /people [get]
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []models.Person
	query := h.DB

	// Фильтрация по полу
	gender := r.URL.Query().Get("gender")
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	// Фильтрация по национальности
	nationality := r.URL.Query().Get("nationality")
	if nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}

	// Получение limit и offset из query-параметров
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10 // Значение по умолчанию
	}

	// Выполнение запроса к базе данных
	query.Limit(limit).Offset(offset).Find(&people)

	// Ответ в формате JSON
	json.NewEncoder(w).Encode(people)
}

// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновляет информацию о человеке по ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param id path int true "ID человека"
// @Param person body models.CreatePersonRequest true "Обновлённые данные"
// @Success 200 {object} models.CreatePersonRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /people/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("id=%s", id)

	var existing models.Person

	// Найти по ID
	if err := h.DB.First(&existing, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление основных полей
	existing.Name = req.Name
	existing.Surname = req.Surname
	existing.Patronymic = req.Patronymic

	// Если имя изменилось — запрашиваем новые значения
	if req.Name != existing.Name {
		age := services.GetAge(req.Name)

		gender := services.GetGender(req.Name)

		nationality := services.GetNationality(req.Name)

		existing.Age = age
		existing.Gender = gender
		existing.Nationality = nationality
	}

	// Сохраняем
	if err := h.DB.Save(&existing).Error; err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existing)
}

// DeletePerson godoc
// @Summary Удалить человека
// @Description Удаляет запись человека по ID
// @Tags people
// @Param id path int true "ID человека"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /people/{id} [delete]
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID параметр из URL
	idStr := chi.URLParam(r, "id")
	log.Printf("Received ID: %s", idStr)

	if idStr == "" {
		log.Printf("No ID provided in the URL!")
		http.Error(w, `{"error":"missing ID"}`, http.StatusBadRequest)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error converting ID: %v", err)
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	var p models.Person
	// Ищем человека по ID
	if err := h.DB.First(&p, id).Error; err != nil {
		log.Printf("Person not found with ID %d", id)
		http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
		return
	}

	// Удаляем человека
	if err := h.DB.Delete(&p).Error; err != nil {
		log.Printf("Error deleting person with ID %d", id)
		http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 (No Content)
	w.WriteHeader(http.StatusNoContent)
}
