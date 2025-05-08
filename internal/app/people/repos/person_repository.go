package repos

import "people/internal/app/people"

type PersonRepository interface {
	Create(person people.Person) (people.Person, error)
}
