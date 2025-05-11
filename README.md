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
Определите интерфейс для базы данных:

В своем репозитории вы можете использовать интерфейс для работы с базой данных. Например:


package repos

import (
    "gorm.io/gorm"
    "context"
)

type DBInterface interface {
    Delete(value interface{}, where ...interface{}) *gorm.DB
}

type GormPersonRepository struct {
    db DBInterface
}

func (r *GormPersonRepository) DeletePerson(ctx context.Context, id int64) error {
    if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
        return err
    }
    return nil
}
Теперь GormPersonRepository использует интерфейс DBInterface, который может быть реализован как для реального *gorm.DB, так и для мока.

Измените мок MockDB для реализации интерфейса DBInterface:


type MockDB struct {
    mock.Mock
}

func (m *MockDB) Delete(value interface{}, where ...interface{}) *gorm.DB {
    args := m.Called(value, where)
    return args.Get(0).(*gorm.DB)
}
Использование моков в тестах:

Теперь вы можете использовать этот интерфейс и моки для тестирования, как это показано ниже:


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
                db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: nil})
            },
            expectedErr: nil,
        },
        {
            name: "person not found",
            id:   2,
            mockDBFunc: func(db *MockDB) {
                db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
            },
            expectedErr: gorm.ErrRecordNotFound,
        },
        {
            name: "database error",
            id:   3,
            mockDBFunc: func(db *MockDB) {
                db.On("Delete", mock.Anything, mock.Anything).Return(&gorm.DB{Error: fmt.Errorf("database error")})
            },
            expectedErr: fmt.Errorf("database error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := new(MockDB)
            repo := &GormPersonRepository{
                db: db,
            }

            tt.mockDBFunc(db)

            err := repo.DeletePerson(context.Background(), tt.id)

            if tt.expectedErr != nil {
                assert.EqualError(t, err, tt.expectedErr.Error())
            } else {
                assert.NoError(t, err)
            }

            db.AssertExpectations(t)
        })
    }
}
Пояснение:
Интерфейс DBInterface: Мы определили интерфейс с методом Delete, который соответствует поведению метода в *gorm.DB. Теперь мы можем использовать и реальный *gorm.DB, и его моки.

Мок MockDB: Мок теперь реализует интерфейс DBInterface, что позволяет использовать его в тестах вместо реального соединения с базой данных.

Репозиторий теперь принимает интерфейс DBInterface, что позволяет вам передавать любые объекты, реализующие этот интерфейс (например, моки или настоящий *gorm.DB).

Теперь ваш тест будет работать без ошибок, так как вы используете интерфейсы, а не конкретные типы.










