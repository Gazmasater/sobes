
DROP TABLE IF EXISTS people;


internal/
└── app/
    └── mydomain/
        ├── usecase/
        │   ├── user_usecase.go        # Бизнес-логика
        │   └── user_usecase_iface.go  # Интерфейс, например UserRepository
        ├── repository/
        │   └── postgres/
        │       └── user_repository.go# Реализация интерфейса
        ├── adapters/
        │   └── http/
        │       └── handler.go         # Использует интерфейс Usecase
        └── domain.go


 curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ivan",
    "surname": "Selivanov",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/person/26"





go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.


git rm --cached textDB


curl -X POST http://localhost:8080/people \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dmitriy",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'

  curl -X DELETE "http://localhost:8080/people/26"


package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id uint) error // Новый метод для удаления
}


package repos

import (
	"context"
	"people/internal/app/people"

	"gorm.io/gorm"
)

// GormPersonRepository реализация PersonRepository через GORM
type GormPersonRepository struct {
	db *gorm.DB
}

// NewPersonRepository создаёт новый GormPersonRepository
func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

// Create сохраняет нового человека в базу данных
func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

// Delete удаляет человека по ID
func (r *GormPersonRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}


package adapterhttp

import (
	"fmt"
	"net/http"
	"strconv"

	"people/internal/app/people/usecase"
)

type Handler struct {
	CreateUC   *usecase.CreatePersonUseCase
	DeleteUC   *usecase.DeletePersonUseCase // Добавляем новый UseCase для удаления
}

func NewHandler(createUC *usecase.CreatePersonUseCase, deleteUC *usecase.DeletePersonUseCase) Handler {
	return Handler{CreateUC: createUC, DeleteUC: deleteUC}
}

func (h Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Обработчик для создания человека
}

func (h Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	idStr := r.URL.Path[len("/persons/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	fmt.Printf("Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.DeleteUC.Execute(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id uint) error {
	return uc.Repo.Delete(ctx, id)
}


gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ go run .
2025/05/09 16:51:04 server started on :8080
Deleting person with ID: 1
DeletePersonUseCase Execute
id=1
2025/05/09 16:51:08 "DELETE http://localhost:8080/people/1 HTTP/1.1" from [::1]:45146 - 000 0B in 753.3µs
2025/05/09 16:51:08 http: panic serving [::1]:45146: runtime error: slice bounds out of range [-1:]
goroutine 21 [running]:
net/http.(*conn).serve.func1()
        /usr/local/go/src/net/http/server.go:1947 +0xbe
panic({0xd6d9a0?, 0xc000246210?})
        /usr/local/go/src/runtime/panic.go:787 +0x132
github.com/go-chi/chi/middleware.prettyStack.decorateFuncCallLine({}, {0xc000b3d270, 0x1e}, 0x1, 0x8)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:130 +0x525
github.com/go-chi/chi/middleware.prettyStack.decorateLine({}, {0xc000b3d270?, 0xdac?}, 0x1, 0x8)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:106 +0x154
github.com/go-chi/chi/middleware.prettyStack.parse({}, {0xc000b3c000, 0xdac, 0xc000019418?}, {0xcb3140, 0x1cf2890})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:89 +0x3b9
github.com/go-chi/chi/middleware.PrintPrettyStack({0xcb3140, 0x1cf2890})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:46 +0x3b
github.com/go-chi/chi/middleware.(*defaultLogEntry).Panic(0x47ff72?, {0xcb3140?, 0x1cf2890?}, {0xc0000194e8?, 0x10000c000019578?, 0x441e00?})
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/logger.go:165 +0x25
github.com/go-chi/chi/middleware.Recoverer.func1.1()
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:28 +0xc8
panic({0xcb3140?, 0x1cf2890?})
        /usr/local/go/src/runtime/panic.go:787 +0x132
people/internal/app/people/usecase.(*DeletePersonUseCase).Execute(0x0, {0xefb8d8, 0xc000b21020}, 0x1)
        /home/gaz358/myprog/sobes/internal/app/people/usecase/usecase.go:75 +0xbb
people/internal/app/people/adapters/adapterhttp.Handler.DeletePerson({0xc00007f2c0?, 0x0?}, {0x7947ba6a1e78, 0xc000127040}, 0xc0001aa780)
        /home/gaz358/myprog/sobes/internal/app/people/adapters/adapterhttp/handlers.go:61 +0x10a
net/http.HandlerFunc.ServeHTTP(0xcda620?, {0x7947ba6a1e78?, 0xc000127040?}, 0xc0001c2dbe?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).routeHTTP(0xc0001105a0, {0x7947ba6a1e78, 0xc000127040}, 0xc0001aa780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:436 +0x1f2
net/http.HandlerFunc.ServeHTTP(0xc000019870?, {0x7947ba6a1e78?, 0xc000127040?}, 0xc000019850?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).ServeHTTP(0xc0001105a0, {0x7947ba6a1e78, 0xc000127040}, 0xc0001aa780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:70 +0x331
github.com/go-chi/chi.(*Mux).Mount.func1({0x7947ba6a1e78, 0xc000127040}, 0xc0001aa780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:311 +0x1c2
net/http.HandlerFunc.ServeHTTP(0xcda620?, {0x7947ba6a1e78?, 0xc000127040?}, 0xc0005046c7?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).routeHTTP(0xc000110540, {0x7947ba6a1e78, 0xc000127040}, 0xc0001aa780)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:436 +0x1f2
net/http.HandlerFunc.ServeHTTP(0x7947ba6a1e78?, {0x7947ba6a1e78?, 0xc000127040?}, 0x1cf2f01?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi/middleware.Recoverer.func1({0x7947ba6a1e78?, 0xc000127040?}, 0xc000b21020?)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/recoverer.go:37 +0x6c
net/http.HandlerFunc.ServeHTTP(0x1cf7480?, {0x7947ba6a1e78?, 0xc000127040?}, 0xc000092aa0?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi/middleware.init.0.RequestLogger.func1.1({0xef9a00, 0xc000b3a000}, 0xc0001aa640)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/middleware/logger.go:57 +0x148
net/http.HandlerFunc.ServeHTTP(0xefb910?, {0xef9a00?, 0xc000b3a000?}, 0x1cf2f60?)
        /usr/local/go/src/net/http/server.go:2294 +0x29
github.com/go-chi/chi.(*Mux).ServeHTTP(0xc000110540, {0xef9a00, 0xc000b3a000}, 0xc0001aa500)
        /home/gaz358/go/pkg/mod/github.com/go-chi/chi@v1.5.5/mux.go:87 +0x2ef
net/http.serverHandler.ServeHTTP({0xc000b20ea0?}, {0xef9a00?, 0xc000b3a000?}, 0x1?)
        /usr/local/go/src/net/http/server.go:3301 +0x8e
net/http.(*conn).serve(0xc000507950, {0xefb8d8, 0xc000b20db0})
        /usr/local/go/src/net/http/server.go:2102 +0x625
created by net/http.(*Server).Serve in goroutine 1
        /usr/local/go/src/net/http/server.go:3454 +0x485
^Csignal: interrupt



