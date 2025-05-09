package usecase

import (
	"context"
	"people/internal/app/people"
	"people/internal/app/people/repos"
)

type CreatePersonUseCase struct {
	Repo repos.PersonRepository
}

type DeletePersonUseCase struct {
	Repo repos.PersonRepository
}

type PersonUseCaseImpl struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
}

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error

	UpdatePerson(ctx context.Context, person people.Person) (people.Person, error)
	GetPersonByID(ctx context.Context, id int64) (people.Person, error)
}

func NewPersonUseCase(
	createUC *CreatePersonUseCase,
	deleteUC *DeletePersonUseCase,
) *PersonUseCaseImpl {
	return &PersonUseCaseImpl{
		CreatePersonUseCase: createUC,
		DeletePersonUseCase: deleteUC,
	}
}

func NewCreatePersonUseCase(repo repos.PersonRepository) *CreatePersonUseCase {
	return &CreatePersonUseCase{Repo: repo}
}

func (uc *CreatePersonUseCase) Execute(ctx context.Context, person people.Person) (people.Person, error) {
	return uc.Repo.Create(ctx, person)
}

func NewDeletePersonUseCase(repo repos.PersonRepository) *DeletePersonUseCase {
	return &DeletePersonUseCase{Repo: repo}
}

func (uc *DeletePersonUseCase) Execute(ctx context.Context, id int64) error {
	return uc.Repo.Delete(ctx, id)
}

func (uc *PersonUseCaseImpl) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCaseImpl) DeletePerson(ctx context.Context, id int64) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}
