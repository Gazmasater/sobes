package repos

import (
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
func (r *GormPersonRepository) Create(person people.Person) (people.Person, error) {
	if err := r.db.Create(&person).Error; err != nil {
		return people.Person{}, err
	}
	return person, nil
}
