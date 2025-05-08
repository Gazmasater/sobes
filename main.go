package main

import (
	"context"
	"net/http"
	"os"
	_ "people/docs"
	"people/pkg/logger"
	"time"

	"people/internal/app/people/adapters/proto_http"
	"people/internal/app/people/repository"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8080
// @BasePath        /
func main() {
	logger.SetLogger(logger.New(zapcore.DebugLevel))

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

	database := repository.Init()
	h := proto_http.Handler{DB: database}

	r := proto_http.SetupRoutes(h)

	logger.Infof(ctx, "Starting server on port: %s", port)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf(ctx, "Server failed: %v", err)
	}
}
