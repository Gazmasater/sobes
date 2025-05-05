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

// CreatePerson godoc
// @Summary      Create a new person
// @Description  Add person by JSON
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        person  body  models.Person  true  "Person"
// @Success      201     {object}  models.Person
// @Failure      400     {object}  map[string]string
// @Router       /people [post]
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
// @Param person body models.Person true "Обновлённые данные"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /people/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Person

	// Поиск существующей записи
	if err := h.DB.First(&p, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	var updated models.Person

	// Декодирование новых данных из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление полей модели
	h.DB.Model(&p).Updates(updated)

	// Ответ с обновлёнными данными
	json.NewEncoder(w).Encode(p)
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
	id := mux.Vars(r)["id"]

	// Удаление записи
	h.DB.Delete(&models.Person{}, id)

	// Ответ без содержимого
	w.WriteHeader(http.StatusNoContent)
}
