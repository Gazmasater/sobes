package usecase

import (
	"context"
	"errors"
	"people/internal/app/people"
	"people/internal/serv"
)

// CreatePersonUseCase структура для обработки создания человека
type CreatePersonUseCase struct {
	PersonRepository repository.PersonRepository
	ExternalService  serv.ExternalService
}

// NewCreatePersonUseCase конструктор для создания нового UseCase
func NewCreatePersonUseCase(pr repository.PersonRepository, es serv.ExternalService) *CreatePersonUseCase {
	return &CreatePersonUseCase{
		PersonRepository: pr,
		ExternalService:  es,
	}
}

// Execute метод, который выполняет логику создания человека
func (uc *CreatePersonUseCase) Execute(ctx context.Context, req people.Person) (people.Person, error) {
	// Получаем данные из внешнего API
	age := uc.ExternalService.GetAge(req.Name)
	gender := uc.ExternalService.GetGender(req.Name)
	nationality := uc.ExternalService.GetNationality(req.Name)

	// Проверяем, что данные валидны
	if age <= 0 || gender == "" || nationality == "" {
		return people.Person{}, errors.New("failed to fetch valid external data")
	}

	// Создаем структуру человека с полученными данными
	person := people.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	// Сохраняем человека в базе данных
	createdPerson, err := uc.PersonRepository.Create(person)
	if err != nil {
		return people.Person{}, err
	}

	// Возвращаем созданного человека
	return createdPerson, nil
}
