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


🔧 Шаг 1: Добавь метод GetPeople в интерфейс PersonUseCase:

type PersonUseCase interface {
	GetPeople(ctx context.Context) ([]people.Person, error)
	// другие методы...
}
🔧 Шаг 2: Имплементация в usecase (если у тебя уже есть personUseCase struct):

type personUseCase struct {
	repo PersonRepository
}

func (uc *personUseCase) GetPeople(ctx context.Context) ([]people.Person, error) {
	return uc.repo.GetPeople(ctx)
}
✅ Шаг 3: Обнови хендлер:

func (h HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	peopleList, err := h.uc.GetPeople(ctx)
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


[{
	"resource": "/home/gaz358/myprog/sobes/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "InvalidIfaceAssign",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "InvalidIfaceAssign"
		}
	},
	"severity": 8,
	"message": "cannot use personUC (variable of type *usecase.PersonUseCaseImpl) as usecase.PersonUseCase value in argument to adapterhttp.NewHandler: *usecase.PersonUseCaseImpl does not implement usecase.PersonUseCase (missing method GetPeople)",
	"source": "compiler",
	"startLineNumber": 85,
	"startColumn": 36,
	"endLineNumber": 85,
	"endColumn": 44
}]














