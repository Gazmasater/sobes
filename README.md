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



package domain

type Person struct {
	ID          uint
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}




package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)


package adapterhttp

import "your_project/internal/app/people/domain"

// Маппинг из domain.Person в PersonResponse.
func toResponse(p domain.Person) PersonResponse {
	return PersonResponse{
		ID:          p.ID,
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         p.Age,
		Gender:      p.Gender,
		Nationality: p.Nationality,
	}
}




// ExternalService интерфейс для внешних сервисов
type ExternalService interface {
	GetAge(name string) int
	GetGender(name string) string
	GetNationality(name string) string
}

// ExternalServiceImpl структура, которая реализует интерфейс ExternalService
type ExternalServiceImpl struct {
	AgifyAPI       string
	GenderizeAPI   string
	NationalizeAPI string
}

// NewExternalService создает новый экземпляр ExternalService с API URL
func NewExternalService() *ExternalServiceImpl {
	return &ExternalServiceImpl{
		AgifyAPI:       os.Getenv("AGIFY_API"),
		GenderizeAPI:   os.Getenv("GENDERIZE_API"),
		NationalizeAPI: os.Getenv("NATIONALIZE_API"),
	}
}

// GetAge получает возраст по имени через API Agify
func (es *ExternalServiceImpl) GetAge(name string) int {
	url := fmt.Sprintf("%s?name=%s", es.AgifyAPI, name)
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0
	}

	return result.Age
}

// GetGender получает пол по имени через API Genderize
func (es *ExternalServiceImpl) GetGender(name string) string {
	url := fmt.Sprintf("%s?name=%s", es.GenderizeAPI, name)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Gender
}

// GetNationality получает национальность по имени через API Nationalize
func (es *ExternalServiceImpl) GetNationality(name string) string {
	url := fmt.Sprintf("%s?name=%s", es.NationalizeAPI, name)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Nationality string `json:"country_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Nationality
}



package usecase

import (
	"context"
	"errors"
	"your_project/internal/app/people/domain"
	"your_project/internal/app/people/repository"
	"your_project/internal/app/people/services"
)

// CreatePersonUseCase структура для обработки создания человека
type CreatePersonUseCase struct {
	PersonRepository repository.PersonRepository
	ExternalService  services.ExternalService
}

// NewCreatePersonUseCase конструктор для создания нового UseCase
func NewCreatePersonUseCase(pr repository.PersonRepository, es services.ExternalService) *CreatePersonUseCase {
	return &CreatePersonUseCase{
		PersonRepository: pr,
		ExternalService:  es,
	}
}

// Execute метод, который выполняет логику создания человека
func (uc *CreatePersonUseCase) Execute(ctx context.Context, req domain.Person) (domain.Person, error) {
	// Получаем данные из внешнего API
	age := uc.ExternalService.GetAge(req.Name)
	gender := uc.ExternalService.GetGender(req.Name)
	nationality := uc.ExternalService.GetNationality(req.Name)

	// Проверяем, что данные валидны
	if age <= 0 || gender == "" || nationality == "" {
		return domain.Person{}, errors.New("failed to fetch valid external data")
	}

	// Создаем структуру человека с полученными данными
	person := domain.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	// Сохраняем человека в базе данных
	createdPerson, err := uc.PersonRepository.Create(person)
	if err != nil {
		return domain.Person{}, err
	}

	// Возвращаем созданного человека
	return createdPerson, nil
}


package main

import (
	"log"
	"net/http"
	"os"
	"your_project/internal/app/people/adapters/adapterhttp"
	"your_project/internal/app/people/repository"
	"your_project/internal/app/people/services"
	"your_project/internal/app/people/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Инициализация репозиториев и сервисов
	db, err := setupDatabase() // Настроить подключение к базе данных
	if err != nil {
		log.Fatal(err)
	}

	// Загружаем ExternalService с переменными окружения
	externalService := services.NewExternalService()

	// Инициализация репозитория
	personRepo := repository.NewPersonRepository(db)

	// Инициализация UseCase
	createPersonUseCase := usecase.NewCreatePersonUseCase(personRepo, externalService)

	// Инициализация хендлеров
	handler := adapterhttp.NewHandler(createPersonUseCase)

	// Настройка маршрутов
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/people", func(r chi.Router) {
		r.Post("/", handler.CreatePerson)
	})

	// Запуск сервера
	http.ListenAndServe(":8080", r)
}



