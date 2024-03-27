package models

import (
	"os"

	"gorm.io/gorm"
)

// type Order struct {
// 	gorm.Model
// 	// ID          uint64 `gorm:"primaryKey" json:"id" binding:"required"`
// 	// ProcID      int    `gorm:"not null" json:"proc_id" binding:"required"`
// 	ProductName string `gorm:"not null" json:"product_name" binding:"required"`
// 	// Quantity    uint64  `gorm:"not null" json:"quantity" binding:"required"`
// 	// Price       float64 `gorm:"not null" json:"price" binding:"required"`
// 	// UserName    string  `gorm:"not null" json:"user_name" binding:"required"`
// 	Timestamp uint64 `gorm:"not null" json:"timestamp" binding:"required"`
// }

type Order struct {
	ProductName string `json:"product_name" binding:"required"`
	Timestamp   uint64 `gorm:"not null" json:"timestamp" binding:"required"`
	ProcID      int    `gorm:"not null" json:"proc_id" binding:"required"`
}

func (o *Order) TableName() string {
	return "orders"
}

// func (o *Order) GetID() uint64 {
// 	return o.ID
// }

func GetById(id uint64, db *gorm.DB) (*Order, error) {
	// get by id
	// return the struct
	var order Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func GetMostRecent(timestamp uint64, db *gorm.DB) (*Order, error) {
	// get by current pid and given timestamp
	// return the struct
	var order Order
	if err := db.Where("proc_id = ? AND timestamp = ?", os.Getpid(), timestamp).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *Order) Insert(db *gorm.DB) error {
	// insert data
	// return the struct
	if err := db.Create(&o).Error; err != nil {
		return err
	}
	return nil
}
