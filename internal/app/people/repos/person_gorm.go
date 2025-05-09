package repos

import (
	"context"
	"fmt"
	"people/internal/app/people"

	"gorm.io/gorm"
)

// GormPersonRepository реализация PersonRepository через GORM
type GormPersonRepository struct {
	db *gorm.DB
}

// NewPersonRepository создаёт новый GormPersonRepository
func NewPersonRepository(db *gorm.DB) *GormPersonRepository {
	return &GormPersonRepository{db: db}
}

// Create сохраняет нового человека в базу данных
func (r *GormPersonRepository) Create(ctx context.Context, person people.Person) (people.Person, error) {
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) FindByID(ctx context.Context, id int64) (people.Person, error) {
	var person people.Person
	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

// Delete удаляет персону из базы данных
func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {

	fmt.Println("Delete")

	if err := r.db.WithContext(ctx).Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormPersonRepository) GetByID(ctx context.Context, id int64) (people.Person, error) {

	var person people.Person

	person.Age = 0
	person.Gender = ""
	person.ID = 0
	person.Name = ""
	person.Nationality = ""
	person.Patronymic = ""
	person.Surname = ""

	fmt.Println("GetByID")

	if err := r.db.WithContext(ctx).First(&person, id).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil

}
