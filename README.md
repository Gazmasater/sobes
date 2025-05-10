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


package repos

import (
	"context"

	"gorm.io/gorm"
	"your_project_path/people" // замените на актуальный путь
)

type GormPersonRepository struct {
	db *gorm.DB
}

func NewGormPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

func (r *GormPersonRepository) CreatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.WithContext(ctx).Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) DeletePerson(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&people.Person{}, id).Error
}

func (r *GormPersonRepository) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	err := r.db.WithContext(ctx).First(&person, id).Error
	return person, err
}

func (r *GormPersonRepository) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	err := r.db.WithContext(ctx).Save(&person).Error
	return person, err
}

func (r *GormPersonRepository) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	var peopleList []people.Person
	query := r.db.WithContext(ctx).Model(&people.Person{})

	// пример простой фильтрации
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Surname != "" {
		query = query.Where("surname ILIKE ?", "%"+filter.Surname+"%")
	}
	if filter.Patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+filter.Patronymic+"%")
	}

	err := query.Find(&peopleList).Error
	return peopleList, err
}

func (r *GormPersonRepository) ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&people.Person{}).
		Where("name = ? AND surname = ? AND patronymic = ?", name, surname, patronymic).
		Count(&count).Error
	return count > 0, err
}





