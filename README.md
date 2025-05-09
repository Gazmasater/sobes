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
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'




1. Интерфейс для создания и удаления персоны:
Создадим интерфейс PersonUseCase в пакете usecase, который будет включать методы для создания и удаления персоны:

go
Копировать
Редактировать
package usecase

import (
	"context"
	"people/internal/app/people"
)

// PersonUseCase интерфейс для работы с персоной
type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error
}
2. Реализация интерфейса в CreatePersonUseCase:
go
Копировать
Редактировать
package usecase

import (
	"context"
	"errors"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

// CreatePersonUseCase структура для обработки создания и удаления человека
type CreatePersonUseCase struct {
	PersonRepository repos.PersonRepository
}

// NewCreatePersonUseCase конструктор для создания нового UseCase
func NewCreatePersonUseCase(pr repos.PersonRepository) *CreatePersonUseCase {
	return &CreatePersonUseCase{
		PersonRepository: pr,
	}
}

// CreatePerson создает новую персону
func (uc *CreatePersonUseCase) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	createdPerson, err := uc.PersonRepository.Create(ctx, req)
	if err != nil {
		return people.Person{}, err
	}
	return createdPerson, nil
}

// DeletePerson удаляет персону по ID
func (uc *CreatePersonUseCase) DeletePerson(ctx context.Context, id int64) error {
	person, err := uc.PersonRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.PersonRepository.Delete(ctx, person)
}
3. Интерфейс для работы с репозиторием:
go
Копировать
Редактировать
package repos

import (
	"context"
	"people/internal/app/people"
)

// PersonRepository интерфейс для работы с репозиторием персоны
type PersonRepository interface {
	// Создание новой персоны
	Create(ctx context.Context, person people.Person) (people.Person, error)
	// Получение персоны по ID
	GetByID(ctx context.Context, id int64) (people.Person, error)
	// Удаление персоны по ID
	Delete(ctx context.Context, person people.Person) error
}
Теперь интерфейс PersonUseCase включает только методы для создания и удаления персоны, как ты и просил.


package repos

import (
	"context"
	"people/internal/app/people"

	"gorm.io/gorm"
)

// GormPersonRepository репозиторий для работы с персоной в базе данных через GORM
type GormPersonRepository struct {
	DB *gorm.DB
}

// NewGormPersonRepository конструктор для создания нового репозитория
func NewGormPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{
		DB: db,
	}
}

// Create создает нового человека в базе данных
func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	// Вставляем нового человека в таблицу
	if err := r.DB.WithContext(ctx).Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

// GetByID получает персону по ID
func (r *GormPersonRepository) GetByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.DB.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

// Delete удаляет персону из базы данных
func (r *GormPersonRepository) Delete(ctx context.Context, person people.Person) error {
	if err := r.DB.WithContext(ctx).Delete(&person).Error; err != nil {
		return err
	}
	return nil
}











