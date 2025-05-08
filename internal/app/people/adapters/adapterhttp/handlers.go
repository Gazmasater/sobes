package adapterhttp

import (
	"encoding/json"
	"net/http"
	"people/internal/services"
	"people/pkg"
	"people/pkg/logger"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
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
// @Param person body CreatePersonRequest true "Данные для создания"
// @Success 200 {object} Person
// @Failure 400 {object} map[string]string
// @Router /people [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(ctx, "Invalid JSON body", "err", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	logger.Debug(ctx, "Parsed request", "request", req)

	// Нормализация
	req.Name = pkg.NormalizeName(req.Name)
	req.Surname = pkg.NormalizeName(req.Surname)
	if req.Patronymic != "" {
		req.Patronymic = pkg.NormalizeName(req.Patronymic)
	}
	logger.Debug(ctx, "Normalized fields", "name", req.Name, "surname", req.Surname, "patronymic", req.Patronymic)

	// Валидация
	if !pkg.IsValidName(req.Name) || !pkg.IsValidName(req.Surname) {
		logger.Warn(ctx, "Validation failed for name/surname", "name", req.Name, "surname", req.Surname)
		http.Error(w, "Name and surname must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}
	if req.Patronymic != "" && !pkg.IsValidName(req.Patronymic) {
		logger.Warn(ctx, "Validation failed for patronymic", "patronymic", req.Patronymic)
		http.Error(w, "Patronymic must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}

	// Внешние API
	age := services.GetAge(req.Name)
	gender := services.GetGender(req.Name)
	nationality := services.GetNationality(req.Name)

	logger.Debug(ctx, "External data fetched", "age", age, "gender", gender, "nationality", nationality)

	p := Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := h.DB.Create(&p).Error; err != nil {
		logger.Error(ctx, "Failed to save person", "err", err)
		http.Error(w, "Failed to save person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "Person created", "id", p.ID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		logger.Error(ctx, "Failed to encode response", "err", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
}

// GetPeople godoc
// @Summary Получить список людей
// @Description Возвращает список людей с возможностью фильтрации, сортировки и пагинации
// @Tags people
// @Accept json
// @Produce json
// @Param name query string false "Имя (например, Ivan)"
// @Param surname query string false "Фамилия (например, Petrov)"
// @Param patronymic query string false "Отчество (например, Ivanovich)"
// @Param age query int false "Возраст (например, 30)"
// @Param gender query string false "Пол (например, male, female)"
// @Param nationality query string false "Национальность (например, Russian, American)"
// @Param sort_by query string false "Поле сортировки (например, name, age, surname)"
// @Param order query string false "Порядок сортировки (asc или desc)"
// @Param limit query int false "Количество возвращаемых записей (по умолчанию 10)"
// @Param offset query int false "Смещение (offset) для пагинации"
// @Success 200 {array} Person
// @Failure 500 {object} map[string]string
// @Router /people [get]
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var people []Person
	query := h.DB

	// Фильтрация
	params := r.URL.Query()

	if gender := params.Get("gender"); gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if nationality := params.Get("nationality"); nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}
	if name := params.Get("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname := params.Get("surname"); surname != "" {
		query = query.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if patronymic := params.Get("patronymic"); patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+patronymic+"%")
	}
	if age := params.Get("age"); age != "" {
		if ageInt, err := strconv.Atoi(age); err == nil {
			query = query.Where("age = ?", ageInt)
		}
	}

	// Сортировка
	sortBy := params.Get("sort_by")
	order := params.Get("order")
	if sortBy != "" {
		if order != "desc" {
			order = "asc"
		}
		// Защита от SQL-инъекций: разрешён только whitelisted список полей
		allowedSorts := map[string]bool{
			"id": true, "name": true, "surname": true,
			"patronymic": true, "age": true, "gender": true, "nationality": true,
		}
		if allowedSorts[sortBy] {
			query = query.Order(sortBy + " " + order)
		}
	}

	// Пагинация
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	if limit == 0 {
		limit = 10
	}
	query = query.Limit(limit).Offset(offset)

	// Получение из базы
	if err := query.Find(&people).Error; err != nil {
		logger.Error(ctx, "DB query failed", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(people); err != nil {
		logger.Error(ctx, "Failed to encode response", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
	}
}

// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновляет информацию о человеке по ID
// @Tags people
// @Accept  json
// @Produce  json
// @Param id path int true "ID человека"
// @Param person body CreatePersonRequest true "Обновлённые данные"
// @Success 200 {object} CreatePersonRequest
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /people/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	logger.Debug(ctx, "Update request received", "id", id)

	var existing Person

	if err := h.DB.First(&existing, id).Error; err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.NotFound(w, r)
		return
	}

	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(ctx, "Invalid JSON body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление полей
	nameChanged := req.Name != existing.Name
	existing.Name = req.Name
	existing.Surname = req.Surname
	existing.Patronymic = req.Patronymic

	if nameChanged {
		logger.Debug(ctx, "Name changed, fetching external data", "old_name", existing.Name, "new_name", req.Name)
		existing.Age = services.GetAge(req.Name)
		existing.Gender = services.GetGender(req.Name)
		existing.Nationality = services.GetNationality(req.Name)
	}

	if err := h.DB.Save(&existing).Error; err != nil {
		logger.Error(ctx, "Failed to update person", "id", id, "err", err)
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "Person updated", "id", id)

	if err := json.NewEncoder(w).Encode(existing); err != nil {
		logger.Error(ctx, "Failed to encode response", "err", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
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

	var p Person
	if err := h.DB.First(&p, id).Error; err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&p).Error; err != nil {
		logger.Error(ctx, "Failed to delete person", "id", id, "err", err)
		http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
		return
	}

	logger.Info(ctx, "Person deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
