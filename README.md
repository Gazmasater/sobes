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


type PersonRepository interface {
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, person people.Person) (people.Person, error)
	GetByID(ctx context.Context, id int64) (people.Person, error)
	ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) // ← вот это обязательно
}

package repos

import (
	"context"
	"fmt"
	"people/internal/app/people"
	"people/pkg"

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

	// Нормализация
	person.Name = pkg.NormalizeName(person.Name)
	person.Surname = pkg.NormalizeName(person.Surname)
	person.Patronymic = pkg.NormalizeName(person.Patronymic)

	// Валидация имени и фамилии (обязательны)
	if !pkg.IsValidName(person.Name) || !pkg.IsValidName(person.Surname) {
		return people.Person{}, fmt.Errorf("invalid name or surname format")
	}

	// Отчество — необязательно, но если есть — проверим
	if len(person.Patronymic) > 0 && !pkg.IsValidName(person.Patronymic) {
		return people.Person{}, fmt.Errorf("invalid patronymic format")
	}

	// Проверка на уникальность
	var existing people.Person
	err := r.db.WithContext(ctx).
		Where("name = ? AND surname = ? AND patronymic = ?", person.Name, person.Surname, person.Patronymic).
		First(&existing).Error

	if err == nil {
		return people.Person{}, fmt.Errorf("person already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return people.Person{}, err
	}

	// Добавление
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {

	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormPersonRepository) GetByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Update(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.WithContext(ctx).Save(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&people.Person{}).
		Where("name = ? AND surname = ? AND patronymic = ?", name, surname, patronymic).
		Count(&count).Error
	return count > 0, err
}











