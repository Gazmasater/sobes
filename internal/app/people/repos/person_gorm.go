package repos

import (
	"context"
	"fmt"
	"people/internal/app/people"
	"people/pkg"

	"gorm.io/gorm"
)

type GormPersonRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

func (r *GormPersonRepository) CreatePerson(ctx context.Context, person people.Person) (people.Person, error) {
	person.Name = pkg.NormalizeName(person.Name)
	person.Surname = pkg.NormalizeName(person.Surname)
	person.Patronymic = pkg.NormalizeName(person.Patronymic)

	if !pkg.IsValidName(person.Name) || !pkg.IsValidName(person.Surname) {
		return people.Person{}, fmt.Errorf("invalid name or surname format")
	}

	if len(person.Patronymic) > 0 && !pkg.IsValidName(person.Patronymic) {
		return people.Person{}, fmt.Errorf("invalid patronymic format")
	}
	var existing people.Person
	err := r.db.WithContext(ctx).
		Where("name = ? AND surname = ? AND patronymic = ?", person.Name, person.Surname, person.Patronymic).
		First(&existing).Error
	if err == nil {
		return people.Person{}, fmt.Errorf("person already exists")
	}

	if err != gorm.ErrRecordNotFound {
		return people.Person{}, err
	}

	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}

	return person, nil
}

func (r *GormPersonRepository) DeletePerson(ctx context.Context, id int64) error {

	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormPersonRepository) UpdatePerson(ctx context.Context, person people.Person) (people.Person, error) {

	if err := r.db.WithContext(ctx).Save(&person).Error; err != nil {
		return people.Person{}, err
	}

	return person, nil
}

func (r *GormPersonRepository) ExistsByFullName(ctx context.Context, name, surname, patronymic string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&people.Person{}).
		Where("name = ? AND surname = ? AND patronymic = ?", name, surname, patronymic).
		Count(&count).Error
	return count > 0, err
}

func (r *GormPersonRepository) GetPersonByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person

	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}

	return person, nil
}

func (r *GormPersonRepository) GetPeople(ctx context.Context, filter people.Filter) ([]people.Person, error) {
	var peopleList []people.Person
	query := r.db.WithContext(ctx)

	if filter.Gender != "" {
		query = query.Where("gender = ?", filter.Gender)
	}

	if filter.Nationality != "" {
		query = query.Where("nationality = ?", filter.Nationality)
	}

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Surname != "" {
		query = query.Where("surname ILIKE ?", "%"+filter.Surname+"%")
	}

	if filter.Patronymic != "" {
		query = query.Where("patronymic ILIKE ?", "%"+filter.Patronymic+"%")
	}

	if filter.Age > 0 {
		query = query.Where("age = ?", filter.Age)
	}

	if filter.SortBy != "" {
		order := "asc"
		if filter.Order == "desc" {
			order = "desc"
		}
		allowed := map[string]bool{
			"id": true, "name": true, "surname": true,
			"patronymic": true, "age": true, "gender": true, "nationality": true,
		}
		if allowed[filter.SortBy] {
			query = query.Order(filter.SortBy + " " + order)
		}
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}

	query = query.Limit(filter.Limit).Offset(filter.Offset)

	if err := query.Find(&peopleList).Error; err != nil {
		return nil, err
	}

	return peopleList, nil
}
