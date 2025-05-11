package mocks

import (
	"context"
	"people/internal/app/people"
)

type MockPersonRepository struct {
	CreateFn    func(ctx context.Context, person people.Person) (people.Person, error)
	DeleteFn    func(ctx context.Context, id int64) error
	GetByIDFn   func(ctx context.Context, id int64) (people.Person, error)
	UpdateFn    func(ctx context.Context, person people.Person) (people.Person, error)
	GetPeopleFn func(ctx context.Context, filter people.Filter) ([]people.Person, error)
}

func (m *MockPersonRepository) CreatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	return m.CreateFn(ctx, person)
}

func (m *MockPersonRepository) DeletePerson(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockPersonRepository) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockPersonRepository) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	return m.UpdateFn(ctx, person)
}

func (m *MockPersonRepository) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	return m.GetPeopleFn(ctx, filter)
}
