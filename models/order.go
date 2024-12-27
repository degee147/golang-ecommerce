package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint    `gorm:"not null"`
	Status string  `gorm:"default:'Pending'"` // Pending, Completed, Cancelled
	Total  float64 `gorm:"not null"`
	Items  []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
}
