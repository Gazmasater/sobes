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
    "surname": "Selivanov",
    "patronymic": "Igorevich"
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



func (es *ExternalServiceImpl) GetGender(ctx context.Context, name string) string {
	url := fmt.Sprintf("%s?name=%s", es.GenderizeAPI, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.Gender
}


type ExternalService interface {
	GetAge(ctx context.Context, name string) int
	GetGender(ctx context.Context, name string) string
	GetNationality(ctx context.Context, name string) string
}


age := uc.ExternalService.GetAge(ctx, person.Name)
gender := uc.ExternalService.GetGender(ctx, person.Name)
nation := uc.ExternalService.GetNationality(ctx, person.Name)


createdPerson, err := h.CreateUC.Execute(r.Context(), person)





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
	"message": "cannot use extService (variable of type *serv.ExternalServiceImpl) as serv.ExternalService value in argument to usecase.NewCreatePersonUseCase: *serv.ExternalServiceImpl does not implement serv.ExternalService (wrong type for method GetAge)\n\t\thave GetAge(context.Context, string) int\n\t\twant GetAge(string) int",
	"source": "compiler",
	"startLineNumber": 41,
	"startColumn": 51,
	"endLineNumber": 41,
	"endColumn": 61
}]







