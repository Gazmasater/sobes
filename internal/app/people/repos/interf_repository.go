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
	Create(ctx context.Context, person people.Person) (people.Person, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, person people.Person) (people.Person, error)
	GetByID(ctx context.Context, id int64) (people.Person, error)
	ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) // ← вот это обязательно
}
