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








package main

import (
	"log"
	"net/http"
	"people/internal/db"
	"people/internal/handlers"
	"people/internal/router"
)

func main() {
	database := db.Init()
	h := handlers.Handler{DB: database}

	r := router.SetupRoutes(h)

	log.Println("API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


package router

import (
	"people/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(h handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/people", func(r chi.Router) {
		r.Post("/", h.CreatePerson)
		r.Get("/", h.GetPeople)
		r.Put("/{id}", h.UpdatePerson)
		r.Delete("/{id}", h.DeletePerson)
	})

	return r
}



✅ 1. Импортировать сгенерированные docs
В файле cmd/server/main.go (или где у тебя точка входа), добавь:


import _ "people/docs" // Путь к пакету с docs, без этого Swagger не заработает
Если у тебя проект в ~/myprog/test, а go.mod начинается с module people, то путь будет корректен.

✅ 2. Добавить маршруты Swagger в Chi
В router/router.go добавь в самый конец:


import (
	httpSwagger "github.com/swaggo/http-swagger"
)

// ...

r.Get("/swagger/*", httpSwagger.WrapHandler)
✅ 3. Пересобери и запусти

go run ./cmd/server
🌐 Swagger доступен по адресу:
http://localhost:8080/swagger/index.html


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










