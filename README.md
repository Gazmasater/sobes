
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


func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {
    // Ищем объект по ID перед удалением
    var person people.Person
    if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
        // Если объект не найден, возвращаем ошибку
        fmt.Println("Error finding person:", err)
        return fmt.Errorf("person not found: %v", err)
    }

    // Печатаем объект, чтобы удостовериться, что он корректно найден
    fmt.Printf("Person found: %+v\n", person)

    // Удаление персоны
    if err := r.db.WithContext(ctx).Delete(&person).Error; err != nil {
        fmt.Println("Error deleting person:", err)
        return fmt.Errorf("failed to delete person: %v", err)
    }

    fmt.Println("Person deleted successfully with ID:", id)
    return nil
}



func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    idStr := chi.URLParam(r, "id")
    logger.Debug(ctx, "Delete request received", "id", idStr)

    if idStr == "" {
        logger.Warn(ctx, "No ID provided in URL")
        http.Error(w, `{"error":"missing ID"}`, http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        logger.Warn(ctx, "Invalid ID format", "id", idStr, "err", err)
        http.Error(w, `{"error":"invalid ID"}`, http.StatusBadRequest)
        return
    }

    fmt.Printf("DeletePerson ID=%d\n", id)

    // Проверяем существование персоны перед удалением
    _, err = h.PersonRepo.GetByID(ctx, id)
    if err != nil {
        logger.Warn(ctx, "Person not found", "id", id, "err", err)
        http.Error(w, `{"error":"person not found"}`, http.StatusNotFound)
        return
    }

    // Удаляем персону
    if err := h.PersonRepo.Delete(ctx, id); err != nil {
        logger.Error(ctx, "Failed to delete person", "id", id, "err", err)
        http.Error(w, fmt.Sprintf(`{"error":"delete failed: %v"}`, err), http.StatusInternalServerError)
        return
    }

    logger.Info(ctx, "Person deleted", "id", id)
    w.WriteHeader(http.StatusNoContent)
}




