package models

import (
	"gorm.io/gorm"
)

// id will be created and incremented by db, else will be provided by user hence have pointers
type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `gorm:"not null" json:"author"`
	Title     *string `josn:"title"`
	Publisher *string `json:"publisher"`
}

// unlike mongodb first you need to create db then start transactions
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
