package models

import "gorm.io/gorm"

// Product represents a product
// @Description Product model for storing products in the system
type Product struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
}
