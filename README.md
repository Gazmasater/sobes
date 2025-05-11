golangci-lint run

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

go test -run=NormalizeName


package yourpackage // замени на название своего пакета

import (
	"testing"
)

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"иван", "Иван"},
		{"  сЕргей", "Сергей"},
		{"ОЛЕГ  ", "Олег"},
		{"", ""},
		{"а", "А"},
		{"   ", ""},
	}

	for _, tt := range tests {
		result := NormalizeName(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeName(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}










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


go test -run=NormalizeName

                          ^
func (r *GormPersonRepository) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	person.Name = pkg.NormalizeName(person.Name)
	person.Surname = pkg.NormalizeName(person.Surname)
	person.Patronymic = pkg.NormalizeName(person.Patronymic)

	if !pkg.IsValidName(person.Name) || !pkg.IsValidName(person.Surname) {
		return people.Person{}, fmt.Errorf("invalid name or surname format")
	}

	if len(person.Patronymic) > 0 && !pkg.IsValidName(person.Patronymic) {
		return people.Person{}, fmt.Errorf("invalid patronymic format")
	}

	// Проверим, не конфликтует ли обновление с уже существующей записью
	var existing people.Person
	err := r.db.WithContext(ctx).
		Where("id <> ? AND name = ? AND surname = ? AND patronymic = ?", person.ID, person.Name, person.Surname, person.Patronymic).
		First(&existing).Error
	if err == nil {
		return people.Person{}, fmt.Errorf("another person with same name already exists")
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return people.Person{}, err
	}

	if err := r.db.WithContext(ctx).Save(&person).Error; err != nil {
		return people.Person{}, err
	}

	return person, nil
}


func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn(ctx, "invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(ctx, "invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		logger.Warn(ctx, "person not found")
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	if req.Name != nil {
		nameChanged := existing.Name != *req.Name
		existing.Name = *req.Name
		if nameChanged {
			existing.Age = h.svc.GetAge(ctx, *req.Name)
			existing.Gender = h.svc.GetGender(ctx, *req.Name)
			existing.Nationality = h.svc.GetNationality(ctx, *req.Name)
		}
	}
	if req.Surname != nil {
		existing.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		existing.Patronymic = *req.Patronymic
	}
	if req.Age != nil {
		existing.Age = *req.Age
	}
	if req.Gender != nil {
		existing.Gender = *req.Gender
	}
	if req.Nationality != nil {
		existing.Nationality = *req.Nationality
	}

	// Обновляем в репозитории (валидация и проверка на дубли внутри)
	updated, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		logger.Warn(ctx, "failed to update person: %v", err)
		http.Error(w, "failed to update person", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ToResponse(updated)); err != nil {
		logger.Error(ctx, "Failed to encode updated response: %v", err)
		http.Error(w, "Failed to encode updated response", http.StatusInternalServerError)
		return
	}
}

