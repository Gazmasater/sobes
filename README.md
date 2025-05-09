
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
    "surname": "Selivanov",
    "patronymic": "Igorevich"
}'

curl -X DELETE "http://localhost:8080/person/26"





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


func main() {
	ctx := logger.ToContext(context.Background(), logger.Global())

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

	dsn := "host=localhost user=postgres password=qwert dbname=people port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	// Миграция таблицы Person
	db.AutoMigrate(&people.Person{})

	// Создание зависимостей
	repo := repos.NewPersonRepository(db)
	extService := serv.NewExternalService()

	// Create and Delete UseCases
	createUC := usecase.NewCreatePersonUseCase(repo)
	deleteUC := usecase.NewDeletePersonUseCase(repo)

	// Объединённый интерфейс
	personUC := usecase.NewPersonUseCase(createUC, deleteUC)

	// Handler принимает один интерфейс
	handler := adapterhttp.NewHandler(personUC)

	// Запуск сервера
	r := adapterhttp.SetupRoutes(handler)
	log.Println("server started on :8080")
	http.ListenAndServe(":8080", r)
}

[{
	"resource": "/home/gaz358/myprog/sobes/main.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "UnusedVar",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UnusedVar"
		}
	},
	"severity": 8,
	"message": "declared and not used: extService",
	"source": "compiler",
	"startLineNumber": 41,
	"startColumn": 2,
	"endLineNumber": 41,
	"endColumn": 12,
	"tags": [
		1
	]
}]

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
	"message": "cannot use repo (variable of type *repos.GormPersonRepository) as repos.PersonRepository value in argument to usecase.NewCreatePersonUseCase: *repos.GormPersonRepository does not implement repos.PersonRepository (wrong type for method Delete)\n\t\thave Delete(context.Context, uint) error\n\t\twant Delete(context.Context, int64) error",
	"source": "compiler",
	"startLineNumber": 44,
	"startColumn": 45,
	"endLineNumber": 44,
	"endColumn": 49
}]




