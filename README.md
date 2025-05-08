people/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   ├── handlers/
│   └── router/
│       └── router.go


go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
go get -u github.com/swaggo/files

Проверь, чтобы $(go env GOPATH)/bin был в $PATH:


export PATH=$PATH:$(go env GOPATH)/bin

2. 📂 Инициализация Swagger
В корне проекта выполни:


swag init
Создастся папка docs с документацией.



func main() {
	// Инициализация логгера
	logger.SetLogger(logger.New(zapcore.DebugLevel))

	// Создаём context
	ctx := logger.ToContext(context.Background(), logger.Global())

	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

	// Порт сервера
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port)

	// Инициализация базы данных
	database := db.Init()

	// Инициализация обработчиков
	h := handlers.Handler{DB: database}

	// Настройка маршрутов
	r := router.SetupRoutes(h)

	logger.Infof(ctx, "Starting server on port: %s", port)

	// Настройка HTTP-сервера с тайм-аутами
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Запуск сервера
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf(ctx, "Server failed: %v", err)
	}
}




