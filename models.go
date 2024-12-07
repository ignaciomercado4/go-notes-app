package main

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
}

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
}
