package people

import (
	"context"
	"people/pkg/logger"

	"gorm.io/gorm"
)

type Person struct {
	ID          uint
	Name        string
	Surname     string
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

type PersonMigration struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"index"`
	Surname    string `gorm:"index"`
	Patronymic string `gorm:"index"`
	Age        int
}

func MigratePersonSchema(ctx context.Context, db *gorm.DB) {
	err := db.AutoMigrate(&PersonMigration{})
	if err != nil {
		logger.Fatalf(ctx, "failed to migrate Person schema: %v", err)
	}
}
