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


go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint --version
golangci-lint run


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
                                 ^
internal/services/services.go:58:10: string `unknown` has 3 occurrences, make it a constant (goconst)
                return "unknown"
                       ^
main.go:39:12: G114: Use of net/http serve function that has no support for setting timeouts (gosec)
        log.Fatal(http.ListenAndServe(":"+port, r))
                  ^
internal/services/services.go:16:23: net/http.Get must not be called (noctx)
        resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
                             ^
internal/services/services.go:35:23: net/http.Get must not be called (noctx)
        resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
                             ^
internal/services/services.go:56:23: net/http.Get must not be called (noctx)
        resp, err := http.Get(fmt.Sprintf("%s?name=%s", apiURL, name))
                             ^
internal/router/router.go:23:2: unnecessary trailing newline (whitespace)
        })
        ^
main.go:21:14: unnecessary leading newline (whitespace)
func main() {
             ^
gaz358@gaz358-BOD-WXX9:~/myprog/test$ 


