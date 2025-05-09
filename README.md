
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
  

func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	fmt.Println("Create")

	var existing people.Person
	err := r.db.WithContext(ctx).Where("name = ? AND surname = ? AND patronymic = ?", person.Name, person.Surname, person.Patronymic).First(&existing).Error
	if err == nil {
		// Такой человек уже есть
		return people.Person{}, fmt.Errorf("person already exists")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return people.Person{}, err
	}

	// Добавляем, если не найден
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/repos/person_gorm.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "cond",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/go/analysis/passes/nilness",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "cond"
		}
	},
	"severity": 4,
	"message": "tautological condition: non-nil != nil",
	"source": "nilness",
	"startLineNumber": 31,
	"startColumn": 10,
	"endLineNumber": 31,
	"endColumn": 10
}]





