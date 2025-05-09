package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	// Создание новой персоны
	Create(ctx context.Context, person people.Person) (people.Person, error)
	// Получение персоны по ID
	GetByID(ctx context.Context, id int64) (people.Person, error)
	// Удаление персоны по ID
	Delete(ctx context.Context, id int64) error
}
