package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// gorm.Open returns two values, capture both
	db, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	// check the error
	if err != nil {
		panic("Failed connecting to database")
	}
	// if no error, assign to the package variable
	DB = db
}
