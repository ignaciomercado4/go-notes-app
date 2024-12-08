package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := connectDatabase()

	r.LoadHTMLGlob("templates/*")

	noteHandler := &NoteHandler{DB: db}

	r.GET("/login", GetLoginForm)

	r.GET("/register", GetRegistrationForm)
	r.POST("/register", noteHandler.CreateUser)

	r.GET("/", noteHandler.GetIndex)

	r.POST("/create", noteHandler.CreateNote)

	r.Run()
}
