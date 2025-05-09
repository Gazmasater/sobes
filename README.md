
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
  



type PersonRepository interface {
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (people.Person, error)        // ✅ для получения по ID
	Update(ctx context.Context, person people.Person) (people.Person, error) // ✅ для обновления
}


type GormPersonRepository struct {
	DB *gorm.DB
}

func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.DB.WithContext(ctx).Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&people.Person{}, id).Error
}

func (r *GormPersonRepository) GetByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.DB.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Update(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.DB.WithContext(ctx).Save(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}


func (uc *PersonUseCaseImpl) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.GetByID(ctx, id)
}

func (uc *PersonUseCaseImpl) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.Update(ctx, person)
}



