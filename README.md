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
	GetPeople(ctx context.Context) ([]people.Person, error)
	ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error)
}


func (r *GormPersonRepository) GetPeople(ctx context.Context) ([]people.Person, error) {
	var peopleList []people.Person
	err := r.db.WithContext(ctx).Find(&peopleList).Error
	if err != nil {
		return nil, err
	}
	return peopleList, nil
}


func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	peopleList, err := h.repo.GetPeople(ctx)
	if err != nil {
		http.Error(w, "failed to get people: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peopleList); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type HTTPHandler struct {
	svc serv.ExternalService

	uc usecase.PersonUseCase
}














