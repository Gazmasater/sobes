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



1. Создание интерфейса PersonRepository
Создайте файл, например, internal/app/repository/person_repository.go:

go
Копировать
Редактировать
package repository

import "people/internal/app/people"

// PersonRepository определяет интерфейс для работы с сущностью Person
type PersonRepository interface {
	Create(person people.Person) (people.Person, error)
}
2. Реализация интерфейса (с GORM, например)
Создайте файл internal/app/repository/person_gorm.go:

go
Копировать
Редактировать
package repository

import (
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
func (r *GormPersonRepository) Create(person people.Person) (people.Person, error) {
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}
3. Настройка GORM и подключение к БД
Пример инициализации GORM в main.go:

go
Копировать
Редактировать
package main

import (
	"log"
	"people/internal/app/repository"
	"people/internal/app/services"
	"people/internal/app/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=qwert dbname=people port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Автоматическая миграция
	db.AutoMigrate(&people.Person{})

	personRepo := repository.NewPersonRepository(db)
	extService := services.NewExternalService() // ваш внешний сервис

	createUC := usecase.NewCreatePersonUseCase(personRepo, extService)

	// Здесь вы подключаете createUC к HTTP-обработчику, как я показывал выше
}



