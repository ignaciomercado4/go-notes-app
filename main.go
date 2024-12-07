package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Note{}, &User{})

	r.LoadHTMLGlob("templates/*")

	noteHandler := &NoteHandler{DB: db}

	r.GET("/", noteHandler.GetIndex)

	r.POST("/create", noteHandler.CreateNote)

	r.Run()
}
