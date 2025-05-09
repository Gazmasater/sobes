package usecase

import (
	"context"
	"errors"
	"fmt"
	"people/internal/app/people"
	"people/internal/app/people/repos"
	"people/internal/serv"
)

type PersonUseCase interface {
	// Создание новой персоны
	CreatePerson(ctx context.Context, req people.Person) (people.Person, error)
	// Удаление персоны по ID
	DeletePerson(ctx context.Context, id int64) error
}

type CreatePersonUseCase struct {
	PersonRepository repos.PersonRepository
	ExternalService  serv.ExternalService
}

func NewCreatePersonUseCase(pr repos.PersonRepository, es serv.ExternalService) *CreatePersonUseCase {
	return &CreatePersonUseCase{
		PersonRepository: pr,
		ExternalService:  es,
	}
}

func (uc *CreatePersonUseCase) Execute(ctx context.Context, req people.Person) (people.Person, error) {

	age := uc.ExternalService.GetAge(ctx, req.Name)
	gender := uc.ExternalService.GetGender(ctx, req.Name)
	nationality := uc.ExternalService.GetNationality(ctx, req.Name)

	fmt.Printf("Execute  AGE=%d\n", age)
	fmt.Printf("Execute  GENDER=%s\n", gender)
	fmt.Printf("Execute  NATION=%s\n", nationality)

	if age <= 0 || gender == "" || nationality == "" {
		return people.Person{}, errors.New("failed to fetch valid external data")
	}

	person := people.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	createdPerson, err := uc.PersonRepository.Create(ctx, person)
	if err != nil {
		return people.Person{}, err
	}

	return createdPerson, nil
}

func (uc *CreatePersonUseCase) DeletePerson(ctx context.Context, id int64) error {
	_, err := uc.PersonRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.PersonRepository.Delete(ctx, id)
}
