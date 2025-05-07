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

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var people []models.Person
	query := h.DB

	// Фильтрация
	params := r.URL.Query()

	if gender := params.Get("gender"); gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if nationality := params.Get("nationality"); nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}
	if name := params.Get("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname := params.Get("surname"); surname != "" {
		query = query.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if patronymic := params.Get("patronymic"); patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+patronymic+"%")
	}
	if age := params.Get("age"); age != "" {
		if ageInt, err := strconv.Atoi(age); err == nil {
			query = query.Where("age = ?", ageInt)
		}
	}

	// Сортировка
	sortBy := params.Get("sort_by")
	order := params.Get("order")
	if sortBy != "" {
		if order != "desc" {
			order = "asc"
		}
		// Защита от SQL-инъекций: разрешён только whitelisted список полей
		allowedSorts := map[string]bool{
			"id": true, "name": true, "surname": true,
			"patronymic": true, "age": true, "gender": true, "nationality": true,
		}
		if allowedSorts[sortBy] {
			query = query.Order(sortBy + " " + order)
		}
	}

	// Пагинация
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	if limit == 0 {
		limit = 10
	}
	query = query.Limit(limit).Offset(offset)

	// Получение из базы
	if err := query.Find(&people).Error; err != nil {
		logger.Error(ctx, "DB query failed", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(people); err != nil {
		logger.Error(ctx, "Failed to encode response", "error", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
	}
}


// GetPeople godoc
// @Summary Получить список людей
// @Description Возвращает список людей с возможностью фильтрации по полу и национальности, а также с пагинацией
// @Tags people
// @Accept json
// @Produce json
// @Param gender query string false "Пол (например, male, female)"
// @Param nationality query string false "Национальность (например, Russian, American)"
// @Param limit query int false "Количество возвращаемых записей (по умолчанию 10)"
// @Param offset query int false "Смещение (offset) для пагинации"
// @Success 200 {array} models.Person
// @Failure 500 {object} map[string]string
// @Router /people [get]



