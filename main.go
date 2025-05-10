package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	//_ "people/docs"
	"people/pkg/logger"
	"time"

	"people/internal/app/people"
	"people/internal/app/people/adapters/adapterhttp"
	"people/internal/app/people/repos"
	"people/internal/app/people/usecase"
	"people/internal/serv"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

// @title           People API
// @version         1.0
// @description     API for managing people.
// @host            localhost:8081
// @BasePath        /
func main() {

	logger.SetLogger(logger.New(zapcore.DebugLevel))

	ctx := logger.ToContext(context.Background(), logger.Global())
	r := chi.NewRouter()

	if err := godotenv.Load(); err != nil {
		logger.Error(ctx, "No .env file found")
	} else {
		logger.Debug(ctx, "Successfully loaded .env file")
	}

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

	port_s := os.Getenv("SERVER_PORT")
	if port_s == "" {
		port_s = "8080"
	}
	logger.Debugf(ctx, "Using port: %s", port_s)

	db.AutoMigrate(&people.Person{})

	repo := repos.NewPersonRepository(db)

	createUC := usecase.NewCreatePersonUseCase(repo)
	deleteUC := usecase.NewDeletePersonUseCase(repo)

	personUC := usecase.NewPersonUseCase(createUC, deleteUC)
	svc := serv.NewExternalService()

	handler := adapterhttp.NewHandler(personUC, svc)

	handler.RegisterRoutes(r)

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
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

}
