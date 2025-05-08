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



gaz358@gaz358-BOD-WXX9:~/myprog/sobes$ golangci-lint run
WARN The linter 'exportloopref' is deprecated (since v1.60.2) due to: Since Go1.22 (loopvar) this linter is no longer relevant. Replaced by copyloopvar. 
ERRO [linters_context] exportloopref: This linter is fully inactivated: it will not produce any reports. 
main.go:48:12: G114: Use of net/http serve function that has no support for setting timeouts (gosec)
        if err := http.ListenAndServe(":"+port, r); err != nil {
                  ^


    func main() {
	// Инициализация логгера
	logger.SetLogger(logger.New(zapcore.DebugLevel))

	// Создаём context
	ctx := logger.ToContext(context.Background(), logger.Global())

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port)

	database := db.Init()
	h := handlers.Handler{DB: database}
	r := router.SetupRoutes(h)

	logger.Infof(ctx, "Starting server on port: %s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatalf(ctx, "Server failed: %v", err)
	}
}






