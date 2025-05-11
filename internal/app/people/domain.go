package people

import (
	"context"
	"people/pkg/logger"

	"gorm.io/gorm"
)

type Person struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"index:idx_name_surname"`
	Surname     string `gorm:"index:idx_name_surname"`
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

type Filter struct {
	Gender      string
	Nationality string
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	SortBy      string
	Order       string
	Limit       int
	Offset      int
}

func MigratePersonSchema(ctx context.Context, db *gorm.DB) {
	err := db.AutoMigrate(&Person{})
	if err != nil {
		logger.Fatalf(ctx, "failed to migrate Person schema: %v", err)
	}

}
