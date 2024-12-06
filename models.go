package main

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	ID      uint
	Title   string
	Content string
}
