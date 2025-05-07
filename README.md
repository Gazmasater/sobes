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


func GetGender(name string) string {
	var res struct {
		Gender string `json:"gender"`
	}

	apiURL := os.Getenv("GENDERIZE_API")
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return ""
	}

	return res.Gender
}

func GetAge(name string) int {
	var res struct {
		Age int `json:"age"`
	}

	apiURL := os.Getenv("AGIFY_API")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", apiURL, name), nil)
	if err != nil {
		// обработка ошибки
		return 0
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0
	}

	return res.Age
}

func GetNationality(name string) string {
	var res struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	apiURL := os.Getenv("NATIONALIZE_API")
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
	if err != nil {
		return unknown
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return unknown
	}

	if len(res.Country) > 0 {
		return res.Country[0].CountryID
	}

	return unknown
}

internal/services/services.go:40:59: Magic number: 5, in <argument> detected (mnd)
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                                                                 ^
internal/services/services.go:20:23: net/http.Get must not be called (noctx)
        resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
                             ^
internal/services/services.go:64:23: net/http.Get must not be called (noctx)
        resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
