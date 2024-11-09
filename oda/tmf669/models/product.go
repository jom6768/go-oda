package models

import (
	"gorm.io/gorm"
)

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2)"`
	CreatedAt   int64   `gorm:"autoCreateTime"`
	UpdatedAt   int64   `gorm:"autoUpdateTime"`
}

// MigrateProduct creates the Product table if it doesn’t exist
func MigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&Product{})
}
