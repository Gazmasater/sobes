package main

import (
	"context"
	"log"
	"net/http"
	_ "people/docs"
	"people/pkg/logger"

	"people/internal/app/people"
	"people/internal/app/people/adapters/adapterhttp"
	"people/internal/app/people/repos"
	"people/internal/app/people/usecase"
	"people/internal/serv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	db.AutoMigrate(&people.Person{})

	repo := repos.NewPersonRepository(db)
	extService := serv.NewExternalService()
	createUC := usecase.NewCreatePersonUseCase(repo, extService)

	deleteUC := usecase.NewDeletePersonUseCase(personRepo)
	handler := adapterhttp.Handler(createUC, deleteUC)

	r := adapterhttp.SetupRoutes(handler)
	log.Println("server started on :8080")
	http.ListenAndServe(":8080", r)
}
