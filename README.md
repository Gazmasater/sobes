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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "people/docs"
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
// @host            localhost:8080
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

package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type CreatePersonUseCase struct {
	Repo repos.PersonRepository
}

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

type PersonUseCaseImpl struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
}

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error

	UpdatePerson(ctx context.Context, person people.Person) (people.Person, error)
	GetPersonByID(ctx context.Context, id int64) (people.Person, error)

	GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error)
}

func NewPersonUseCase(
	createUC *CreatePersonUseCase,
	deleteUC *DeletePersonUseCase,
) *PersonUseCaseImpl {
	return &PersonUseCaseImpl{
		CreatePersonUseCase: createUC,
		DeletePersonUseCase: deleteUC,
	}
}

func NewCreatePersonUseCase(repo repos.PersonRepository) *CreatePersonUseCase {
	return &CreatePersonUseCase{Repo: repo}
}

func (uc *CreatePersonUseCase) Execute(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.Repo.CreatePerson(ctx, person)
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id int64) error {
	return uc.Repo.DeletePerson(ctx, id)
}

func (uc *PersonUseCaseImpl) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCaseImpl) DeletePerson(ctx context.Context, id int64) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}

func (uc *PersonUseCaseImpl) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.GetPersonByID(ctx, id)
}

func (uc *PersonUseCaseImpl) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Repo.UpdatePerson(ctx, person)
}

func (uc *PersonUseCaseImpl) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	return uc.CreatePersonUseCase.Repo.GetPeople(ctx, filter)
}

package repos

import (
	"context"
	"fmt"
	"people/internal/app/people"
	"people/pkg"

	"gorm.io/gorm"
)

// GormPersonRepository реализация PersonRepository через GORM
type GormPersonRepository struct {
	db *gorm.DB
}

// NewPersonRepository создаёт новый GormPersonRepository
func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

// Create сохраняет нового человека в базу данных
func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {

	// Нормализация
	person.Name = pkg.NormalizeName(person.Name)
	person.Surname = pkg.NormalizeName(person.Surname)
	person.Patronymic = pkg.NormalizeName(person.Patronymic)

	// Валидация имени и фамилии (обязательны)
	if !pkg.IsValidName(person.Name) || !pkg.IsValidName(person.Surname) {
		return people.Person{}, fmt.Errorf("invalid name or surname format")
	}

	// Отчество — необязательно, но если есть — проверим
	if len(person.Patronymic) > 0 && !pkg.IsValidName(person.Patronymic) {
		return people.Person{}, fmt.Errorf("invalid patronymic format")
	}

	// Проверка на уникальность
	var existing people.Person
	err := r.db.WithContext(ctx).
		Where("name = ? AND surname = ? AND patronymic = ?", person.Name, person.Surname, person.Patronymic).
		First(&existing).Error

	if err == nil {
		return people.Person{}, fmt.Errorf("person already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return people.Person{}, err
	}

	// Добавление
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {

	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormPersonRepository) GetByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Update(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.WithContext(ctx).Save(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&people.Person{}).
		Where("name = ? AND surname = ? AND patronymic = ?", name, surname, patronymic).
		Count(&count).Error
	return count > 0, err
}

func (r *GormPersonRepository) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	var peopleList []people.Person
	query := r.db.WithContext(ctx)

	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}
	if filter.Nationality != "" {
		query = query.Where("nationality = ?", filter.Nationality)
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Surname != "" {
		query = query.Where("surname ILIKE ?", "%"+filter.Surname+"%")
	}
	if filter.Patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+filter.Patronymic+"%")
	}
	if filter.Age > 0 {
		query = query.Where("age = ?", filter.Age)
	}

	if filter.SortBy != "" {
		order := "asc"
		if filter.Order == "desc" {
			order = "desc"
		}
		allowed := map[string]bool{
			"id": true, "name": true, "surname": true,
			"patronymic": true, "age": true, "gender": true, "nationality": true,
		}
		if allowed[filter.SortBy] {
			query = query.Order(filter.SortBy + " " + order)
		}
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	query = query.Limit(filter.Limit).Offset(filter.Offset)

	if err := query.Find(&peopleList).Error; err != nil {
		return nil, err
	}
	return peopleList, nil
}

package repos

import (
	"context"
	"people/internal/app/people"
)

// type PersonRepository interface {
// 	Create(ctx context.Context, person people.Person) (people.Person, error)
// 	Delete(ctx context.Context, id int64) error
// }

type PersonRepository interface {
	CreatePerson(ctx context.Context, person people.Person) (people.Person, error)
	DeletePerson(ctx context.Context, id int64) error
	UpdatePerson(ctx context.Context, person people.Person) (people.Person, error)
	GetPersonByID(ctx context.Context, id int64) (people.Person, error)
	GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error)
}

package adapterhttp

import (
	"encoding/json"
	"net/http"
	"people/internal/app/people"
	"people/internal/app/people/usecase"
	"people/internal/serv"
	"people/pkg/logger"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type HTTPHandler_interf interface {
	RegisterRoutes(r chi.Router)
}

type HTTPHandler struct {
	svc serv.ExternalService

	uc usecase.PersonUseCase
}

func NewHandler(uc usecase.PersonUseCase, svc serv.ExternalService) HTTPHandler_interf {
	return &HTTPHandler{uc: uc, svc: svc}
}

// CreatePerson godoc
// @Summary      Create person
// @Description  Creates a new person with enriched data
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        person  body      CreatePersonRequest  true  "Person to create"
// @Success      200     {object}  PersonResponse
// @Failure      400     {string}  string  "invalid request body"
// @Failure      500     {string}  string  "internal server error"
// @Router       /people [post]
func (h HTTPHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(ctx, "invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем доп. данные
	age := h.svc.GetAge(ctx, req.Name)
	gender := h.svc.GetGender(ctx, req.Name)
	nationality := h.svc.GetNationality(ctx, req.Name)
	logger.Debug(ctx, "CreatePerson AGE=%d GENDER=%s NATION=%s", age, gender, nationality)

	person := people.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdPerson, err := h.uc.CreatePerson(ctx, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ToResponse(createdPerson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeletePerson godoc
// @Summary      Delete person
// @Description  Deletes person by ID
// @Tags         people
// @Produce      json
// @Param        id   path      int64  true  "Person ID"
// @Success      204  {string}  string  "no content"
// @Failure      400  {string}  string  "invalid id"
// @Failure      500  {string}  string  "internal server error"
// @Router       /people/{id} [delete]
func (h HTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Path[len("/people/"):]

	logger.Debug(ctx, "DeletePerson URL=%s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.Debug(ctx, "Deleting person with ID: %d\n", id)

	// Вызываем UseCase для удаления
	err = h.uc.DeletePerson(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePerson godoc
// @Summary      Update person
// @Description  Updates person by ID and enriches if name changed
// @Tags         people
// @Accept       json
// @Produce      json
// @Param        id      path      int64                 true  "Person ID"
// @Param        person  body      UpdatePersonRequest   true  "Updated person (partial)"
// @Success      200     {object}  PersonResponse
// @Failure      400     {string}  string  "invalid request body or id"
// @Failure      404     {string}  string  "person not found"
// @Failure      500     {string}  string  "failed to update person"
// @Router       /people/{id} [put]
func (h HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем ID из URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn(ctx, "invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Декодируем JSON-тело запроса
	var req UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn(ctx, "invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем существующего человека
	existing, err := h.uc.GetPersonByID(ctx, id)
	if err != nil {
		logger.Warn(ctx, "person not found")
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	// Проверяем и обновляем только переданные поля
	if req.Name != nil {
		nameChanged := existing.Name != *req.Name
		existing.Name = *req.Name
		if nameChanged {
			existing.Age = h.svc.GetAge(ctx, *req.Name)
			existing.Gender = h.svc.GetGender(ctx, *req.Name)
			existing.Nationality = h.svc.GetNationality(ctx, *req.Name)
		}
	}
	if req.Surname != nil {
		existing.Surname = *req.Surname
	}
	if req.Patronymic != nil {
		existing.Patronymic = *req.Patronymic
	}
	if req.Age != nil {
		existing.Age = *req.Age
	}
	if req.Gender != nil {
		existing.Gender = *req.Gender
	}
	if req.Nationality != nil {
		existing.Nationality = *req.Nationality
	}

	// Обновляем запись
	updated, err := h.uc.UpdatePerson(ctx, existing)
	if err != nil {
		logger.Warn(ctx, "failed to update person")
		http.Error(w, "failed to update person", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(updated))
}

func (h *HTTPHandler) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/people", h.CreatePerson)
	r.Delete("/people/{id}", h.DeletePerson)
	r.Put("/people/{id}", h.UpdatePerson)
	r.Get("/people", h.GetPeople)

}

func (h *HTTPHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()

	age, _ := strconv.Atoi(params.Get("age"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	if limit == 0 {
		limit = 10
	}

	filter := people.Filter{
		Gender:      params.Get("gender"),
		Nationality: params.Get("nationality"),
		Name:        params.Get("name"),
		Surname:     params.Get("surname"),
		Patronymic:  params.Get("patronymic"),
		Age:         age,
		SortBy:      params.Get("sort_by"),
		Order:       params.Get("order"),
		Limit:       limit,
		Offset:      offset,
	}

	peopleList, err := h.uc.GetPeople(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to get people", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peopleList); err != nil {
		logger.Error(ctx, "failed to encode response", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
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
	"message": "cannot use repo (variable of type *repos.GormPersonRepository) as repos.PersonRepository value in argument to usecase.NewCreatePersonUseCase: *repos.GormPersonRepository does not implement repos.PersonRepository (missing method CreatePerson)",
	"source": "compiler",
	"startLineNumber": 79,
	"startColumn": 45,
	"endLineNumber": 79,
	"endColumn": 49
}]




