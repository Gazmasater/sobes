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
	fmt.Println("Create")

	var existing people.Person
	err := r.db.WithContext(ctx).Where("name = ? AND surname = ? AND patronymic = ?", person.Name, person.Surname, person.Patronymic).First(&existing).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return people.Person{}, err
		}
	} else {
		// Такой человек уже есть
		return people.Person{}, fmt.Errorf("person already exists")
	}

	// Добавляем, если не найден
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}

func (r *GormPersonRepository) Delete(ctx context.Context, id int64) error {

	fmt.Println("Delete")

	if err := r.db.Delete(&people.Person{}, id).Error; err != nil {
		return err
	}
	return nil
}
