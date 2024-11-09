package models

import (
	"gorm.io/gorm"
)

type Service struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	CreatedAt   int64  `gorm:"autoCreateTime"`
	UpdatedAt   int64  `gorm:"autoUpdateTime"`
}

// MigrateService creates the Service table if it doesn’t exist
func MigrateService(db *gorm.DB) error {
	return db.AutoMigrate(&Service{})
}
