package repos

import "people/internal/app/people"

// PersonRepository определяет интерфейс для работы с сущностью Person
type PersonRepository interface {
	Create(person people.Person) (people.Person, error)
}
