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

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Нормализация (первая буква заглавная, остальные строчные)
	req.Name = normalizeName(req.Name)
	req.Surname = normalizeName(req.Surname)
	if req.Patronymic != "" {
		req.Patronymic = normalizeName(req.Patronymic)
	}

	// Валидация имени и фамилии (обязательные), отчество — опционально
	if !isValidName(req.Name) || !isValidName(req.Surname) {
		http.Error(w, "Name and surname must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}
	if req.Patronymic != "" && !isValidName(req.Patronymic) {
		http.Error(w, "Patronymic must contain only letters and start with a capital letter", http.StatusBadRequest)
		return
	}

	// Получение данных из внешних API
	age, err := services.GetAge(req.Name)
	if err != nil {
		http.Error(w, "Failed to fetch age: "+err.Error(), http.StatusInternalServerError)
		return
	}

	gender, err := services.GetGender(req.Name)
	if err != nil {
		http.Error(w, "Failed to fetch gender: "+err.Error(), http.StatusInternalServerError)
		return
	}

	nationality, err := services.GetNationality(req.Name)
	if err != nil {
		http.Error(w, "Failed to fetch nationality: "+err.Error(), http.StatusInternalServerError)
		return
	}

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



func normalizeName(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func isValidName(name string) bool {
	match, _ := regexp.MatchString(`^[\p{L}]+$`, name) // Только буквы (любой алфавит)
	if !match {
		return false
	}
	runes := []rune(name)
	return unicode.IsUpper(runes[0])
}










