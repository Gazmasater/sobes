package adapterhttp

import "people/internal/app/people"

func ToResponse(p people.Person) PersonResponse {
	return PersonResponse{
		ID:          p.ID,
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         p.Age,
		Gender:      p.Gender,
		Nationality: p.Nationality,
	}
}
