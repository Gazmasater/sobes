
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

  curl -X DELETE "http://localhost:8080/people/26"
  



func (uc *PersonUseCaseImpl) UpdatePerson(ctx context.Context, updated people.Person) (people.Person, error) {
	existing, err := uc.CreatePersonUseCase.Repo.GetByID(ctx, updated.ID)
	if err != nil {
		return people.Person{}, fmt.Errorf("person not found: %w", err)
	}

	nameChanged := existing.Name != updated.Name

	// Обновляем все поля
	existing.Name = updated.Name
	existing.Surname = updated.Surname
	existing.Patronymic = updated.Patronymic
	existing.Age = updated.Age
	existing.Gender = updated.Gender
	existing.Nationality = updated.Nationality

	if nameChanged {
		existing.Age = uc.CreatePersonUseCase.ExtSvc.GetAge(ctx, updated.Name)
		existing.Gender = uc.CreatePersonUseCase.ExtSvc.GetGender(ctx, updated.Name)
		existing.Nationality = uc.CreatePersonUseCase.ExtSvc.GetNationality(ctx, updated.Name)
	}

	return uc.CreatePersonUseCase.Repo.Update(ctx, existing)
}


