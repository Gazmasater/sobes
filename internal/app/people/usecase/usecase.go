package usecase

import (
	"context"
	"people/internal/app/people"
)

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error
}

type PersonUseCaseImpl struct {
	CreatePersonUseCase *CreatePersonUseCase
	DeletePersonUseCase *DeletePersonUseCase
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

func (uc *PersonUseCaseImpl) CreatePerson(ctx context.Context, req people.Person) (people.Person, error) {
	return uc.CreatePersonUseCase.Execute(ctx, req)
}

func (uc *PersonUseCaseImpl) DeletePerson(ctx context.Context, id int64) error {
	return uc.DeletePersonUseCase.Execute(ctx, id)
}
