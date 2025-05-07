people/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   ├── handlers/
│   └── router/
│       └── router.go


go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.

package main

import (
	"log"
	"net/http"

	"people/docs" // путь к swagger docs
	"people/internal/db"
	"people/internal/handlers"
	"people/internal/router"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /

func main() {
	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	// Swagger endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


4. ✏️ Пример аннотации в обработчике
📁 internal/handlers/person.go

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
	// ...
}
5. 🔁 Обнови Swagger при изменении кода
Каждый раз после изменения аннотаций:


swag init
✅ Swagger будет доступен по адресу:

http://localhost:8080/swagger/index.html


gaz358@gaz358-BOD-WXX9:~/myprog/test$ golangci-lint run
WARN The linter 'exportloopref' is deprecated (since v1.60.2) due to: Since Go1.22 (loopvar) this linter is no longer relevant. Replaced by copyloopvar. 
ERRO [linters_context] exportloopref: This linter is fully inactivated: it will not produce any reports. 
internal/handlers/handlers.go:80:27: Error return value of `(*encoding/json.Encoder).Encode` is not checked (errcheck)
        json.NewEncoder(w).Encode(p)
                                 ^
internal/handlers/handlers.go:123:27: Error return value of `(*encoding/json.Encoder).Encode` is not checked (errcheck)
        json.NewEncoder(w).Encode(people)
                                 ^
internal/handlers/handlers.go:180:27: Error return value of `(*encoding/json.Encoder).Encode` is not checked (errcheck)
        json.NewEncoder(w).Encode(existing)


 package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"people/internal/models"
	"people/internal/pkg"
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
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Нормализация (первая буква заглавная, остальные строчные)
	req.Name = pkg.NormalizeName(req.Name)
	req.Surname = pkg.NormalizeName(req.Surname)
	if req.Patronymic != "" {
		req.Patronymic = pkg.NormalizeName(req.Patronymic)
	}

	// Валидация имени и фамилии (обязательные), отчество — опционально
	if !pkg.IsValidName(req.Name) || !pkg.IsValidName(req.Surname) {
		http.Error(w, "Name and surname must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}
	if req.Patronymic != "" && !pkg.IsValidName(req.Patronymic) {
		http.Error(w, "Patronymic must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}

	// Получение данных из внешних API
	age := services.GetAge(req.Name)

	gender := services.GetGender(req.Name)

	nationality := services.GetNationality(req.Name)

	// Создание и сохранение записи
	p := models.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	if err := h.DB.Create(&p).Error; err != nil {
		http.Error(w, "Failed to save person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ
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



