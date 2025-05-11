package repos

import (
	"context"
	"people/internal/app/people"
)

type PersonRepository interface {
	CreatePerson(ctx context.Context, person people.Person) (people.Person, error)
	DeletePerson(ctx context.Context, id int64) error
	GetPersonByID(ctx context.Context, id int64) (people.Person, error)
	UpdatePerson(ctx context.Context, person people.Person) (people.Person, error)
	GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error)
}
