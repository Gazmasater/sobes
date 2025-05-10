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

  gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ go run .
{"lvl":"debug","ts":"2025-05-10T22:02:20.961+0300","msg":"Successfully loaded .env file"}
{"lvl":"debug","ts":"2025-05-10T22:02:20.965+0300","msg":"Using port: 8080"}
{"lvl":"info","ts":"2025-05-10T22:02:21.008+0300","msg":"Starting server on port: 8080"}
2025/05/10 22:02:42 "GET http://localhost:8080/people/3 HTTP/1.1" from [::1]:59528 - 405 0B in 47.854µs








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


package people

import "gorm.io/gorm"

// Person структура для представления человека
type Person struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
	Name       string `json:"name" gorm:"index"`
	Surname    string `json:"surname" gorm:"index"`
	Patronymic string `json:"patronymic" gorm:"index"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Nationality string `json:"nationality"`
}

// AddIndexes добавляет индексы для оптимизации запросов
func (Person) AddIndexes(db *gorm.DB) error {
	if err := db.Model(&Person{}).AddIndex("idx_name", "name").Error; err != nil {
		return err
	}
	if err := db.Model(&Person{}).AddIndex("idx_surname", "surname").Error; err != nil {
		return err
	}
	if err := db.Model(&Person{}).AddIndex("idx_patronymic", "patronymic").Error; err != nil {
		return err
	}
	return nil
}



func main() {
	// Инициализация логгера и контекста
	logger.SetLogger(logger.New(zapcore.DebugLevel))
	ctx := logger.ToContext(context.Background(), logger.Global())

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

	// Подключение к базе данных
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	// Автоматическая миграция для таблицы Person
	db.AutoMigrate(&people.Person{})

	// Создание индексов для оптимизации поиска
	if err := people.Person{}.AddIndexes(db); err != nil {
		logger.Fatal(ctx, "Failed to create indexes", "error", err)
	}

	// Инициализация репозитория, юзкейсов и обработчиков
	repo := repos.NewPersonRepository(db)
	createUC := usecase.NewCreatePersonUseCase(repo)
	deleteUC := usecase.NewDeletePersonUseCase(repo)
	personUC := usecase.NewPersonUseCase(createUC, deleteUC)
	svc := serv.NewExternalService()
	handler := adapterhttp.NewHandler(personUC, svc)

	// Регистрируем маршруты и запускаем сервер
	handler.RegisterRoutes(r)

	port_s := os.Getenv("SERVER_PORT")
	if port_s == "" {
		port_s = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port_s)

	logger.Infof(ctx, "Starting server on port: %s", port_s)

	srv := &http.Server{
		Addr:         ":" + port_s,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf(ctx, "Server failed: %v", err)
	}
}


[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/domain.go",
	"owner": "_generated_diagnostic_collection_name_#0",
	"code": {
		"value": "MissingFieldOrMethod",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "MissingFieldOrMethod"
		}
	},
	"severity": 8,
	"message": "db.Model(&Person{}).AddIndex undefined (type *gorm.DB has no field or method AddIndex)",
	"source": "compiler",
	"startLineNumber": 29,
	"startColumn": 32,
	"endLineNumber": 29,
	"endColumn": 40
}]






