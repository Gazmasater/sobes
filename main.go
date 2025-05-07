package main

import (
	"context"
	"net/http"
	"os"
	_ "people/docs"

	"people/internal/db"
	"people/internal/handlers"
	"people/internal/pkg/logger"
	"people/internal/router"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /
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
