package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := connectDatabase()

	r.LoadHTMLGlob("templates/*")

	noteHandler := &NoteHandler{DB: db}

	r.GET("/", noteHandler.GetIndex)

	r.POST("/create", noteHandler.CreateNote)

	r.Run()
}
