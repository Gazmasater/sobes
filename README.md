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

	if err := json.NewEncoder(w).Encode(existing); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
}

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


