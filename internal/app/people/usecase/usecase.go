package usecase

import (
	"context"
	"errors"
	"fmt"
	"people/internal/app/people"
	"people/internal/app/people/repos"
	"people/internal/serv"
)

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

	fmt.Printf("Execute   NAME=%s\n", req.Name)
	age := uc.ExternalService.GetAge(req.Name)
	fmt.Printf("Execute AGE=%d\n", age)
	gender := uc.ExternalService.GetGender(req.Name)
	fmt.Printf("Execute GENDER=%s\n", gender)

	nationality := uc.ExternalService.GetNationality(req.Name)
	fmt.Printf("Execute NATIONAL=%s\n", nationality)

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

	createdPerson, err := uc.PersonRepository.Create(person)
	if err != nil {
		return people.Person{}, err
	}

	return createdPerson, nil
}
