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
package repos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(value, where)
	return args.Get(0).(*gorm.DB)
}

func TestDeletePerson(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		mockDBFunc  func(db *MockDB)
		expectedErr error
	}{
		{
			name: "successful deletion",
			id:   1,
			mockDBFunc: func(db *MockDB) {
				// Мокаем успешное удаление
				db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: nil})
			},
			expectedErr: nil,
		},
		{
			name: "person not found",
			id:   2,
			mockDBFunc: func(db *MockDB) {
				// Мокаем ошибку при удалении (например, если записи нет)
				db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "database error",
			id:   3,
			mockDBFunc: func(db *MockDB) {
				// Мокаем ошибку базы данных
				db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: fmt.Errorf("database error")})
			},
			expectedErr: fmt.Errorf("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация мока базы данных и репозитория
			db := new(MockDB)
			repo := &GormPersonRepository{
				db: db, // Подключаем мок базы данных к репозиторию
			}

			// Настройка мока для текущего теста
			tt.mockDBFunc(db)

			// Запуск тестируемого метода
			err := repo.DeletePerson(context.Background(), tt.id)

			// Сравнение ошибки с ожидаемой
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверка, что методы мока были вызваны
			db.AssertExpectations(t)
		})
	}
}

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/repos/del_test.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "IncompatibleAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "IncompatibleAssign"
		}
	},
	"severity": 8,
	"message": "cannot use db (variable of type *MockDB) as *gorm.DB value in struct literal",
	"source": "compiler",
	"startLineNumber": 63,
	"startColumn": 9,
	"endLineNumber": 63,
	"endColumn": 11
}]

