package main

import (
	"go-notes-app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Note{})

	return db
}
