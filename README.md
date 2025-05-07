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
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	logger.Debug(ctx, "Update request received", "id", id)

	var existing models.Person

	if err := h.DB.First(&existing, id).Error; err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.NotFound(w, r)
		return
	}

	var req models.CreatePersonRequest
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

	var p models.Person
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


