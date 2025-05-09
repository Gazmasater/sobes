package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id int64) error // Новый метод для удаления

}
