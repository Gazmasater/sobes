go install github.com/swaggo/swag/cmd/swag@latest

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
    "surname": "Seli",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/people/1"


curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович",
    "age": 35,
    "gender": "male",
    "nationality": "russian"
  }'


  curl -X GET http://localhost:8080/people

  gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ go run .
{"lvl":"debug","ts":"2025-05-10T22:02:20.961+0300","msg":"Successfully loaded .env file"}
{"lvl":"debug","ts":"2025-05-10T22:02:20.965+0300","msg":"Using port: 8080"}
{"lvl":"info","ts":"2025-05-10T22:02:21.008+0300","msg":"Starting server on port: 8080"}
2025/05/10 22:02:42 "GET http://localhost:8080/people/3 HTTP/1.1" from [::1]:59528 - 405 0B in 47.854µs








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

  curl -X DELETE "http://localhost:8080/people/5"


  curl -X PUT http://localhost:8080/people/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alexey",
    "surname": "Ivanov",
    "patronymic": "Sergeevich",
    "age": 30,
    "gender": "male",
    "nationality": "ru"
  }'

  
swag init -g cmd/main.go -o docs




package people

type Person struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type Filter struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         *int
	Gender      string
	Nationality string
	Limit       int
	Offset      int
	SortBy      string
	Order       string
}


package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error)
}


package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type GetPeopleUseCase struct {
	Repo repos.PersonRepository
}

func NewGetPeopleUseCase(repo repos.PersonRepository) *GetPeopleUseCase {
	return &GetPeopleUseCase{Repo: repo}
}

func (uc *GetPeopleUseCase) Execute(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	return uc.Repo.GetPeople(ctx, filter)
}


package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"people/internal/app/people"
	"people/internal/app/people/usecase"
)

type Handler struct {
	GetPeopleUC *usecase.GetPeopleUseCase
}

func NewHandler(getPeopleUC *usecase.GetPeopleUseCase) *Handler {
	return &Handler{
		GetPeopleUC: getPeopleUC,
	}
}

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()

	filter := people.Filter{
		Name:        params.Get("name"),
		Surname:     params.Get("surname"),
		Patronymic:  params.Get("patronymic"),
		Gender:      params.Get("gender"),
		Nationality: params.Get("nationality"),
		SortBy:      params.Get("sort_by"),
		Order:       params.Get("order"),
		Limit:       parseInt(params.Get("limit"), 10),
		Offset:      parseInt(params.Get("offset"), 0),
	}

	if ageStr := params.Get("age"); ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			filter.Age = &age
		}
	}

	peopleList, err := h.GetPeopleUC.Execute(ctx, filter)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(peopleList)
}

func parseInt(val string, def int) int {
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}


package repos

import (
	"context"
	"people/internal/app/people"

	"gorm.io/gorm"
)

type GormPersonRepository struct {
	DB *gorm.DB
}

func NewGormPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{DB: db}
}

func (r *GormPersonRepository) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	query := r.DB.WithContext(ctx).Model(&people.Person{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Surname != "" {
		query = query.Where("surname ILIKE ?", "%"+filter.Surname+"%")
	}
	if filter.Patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+filter.Patronymic+"%")
	}
	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}
	if filter.Nationality != "" {
		query = query.Where("nationality = ?", filter.Nationality)
	}
	if filter.Age != nil {
		query = query.Where("age = ?", *filter.Age)
	}

	if filter.SortBy != "" {
		order := "ASC"
		if filter.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(filter.SortBy + " " + order)
	}

	var peopleList []people.Person
	err := query.Limit(filter.Limit).Offset(filter.Offset).Find(&peopleList).Error
	return peopleList, err
}



