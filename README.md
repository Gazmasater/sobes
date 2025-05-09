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

curl -X DELETE "http://localhost:8080/person/123"





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
    "name": "Иван",
    "surname": "Иванов",
    "patronymic": "Иванович"
  }'


	createUC := usecase.NewCreatePersonUseCase(repo, extService)


func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	logger.Debug(ctx, "Delete request received", "id", idStr)

	// Проверка на пустой ID в URL
	if idStr == "" {
		logger.Warn(ctx, "No ID provided in URL")
		http.Error(w, `{"error":"missing ID"}`, http.StatusBadRequest)
		return
	}

	// Преобразование строки ID в целое число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Warn(ctx, "Invalid ID format", "id", idStr, "err", err)
		http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Проверка существования персоны в базе данных
	var p Person
	if err := h.DB.First(&p, id).Error; err != nil {
		logger.Warn(ctx, "Person not found", "id", id, "err", err)
		http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
		return
	}

	// Удаление персоны из базы данных
	if err := h.DB.Delete(&p).Error; err != nil {
		logger.Error(ctx, "Failed to delete person", "id", id, "err", err)
		http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
		return
	}

	// Успешное удаление
	logger.Info(ctx, "Person deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

[{
	"resource": "/home/gaz358/myprog/sobes/internal/app/people/adapters/adapterhttp/handlers.go",
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
	"message": "h.DB undefined (type *Handler has no field or method DB)",
	"source": "compiler",
	"startLineNumber": 71,
	"startColumn": 14,
	"endLineNumber": 71,
	"endColumn": 16
}]
