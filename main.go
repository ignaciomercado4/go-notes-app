package main

import (
	"go-notes-app/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.ConnectDatabase()

	r.LoadHTMLGlob("templates/*")

	noteHandler := &NoteHandler{DB: db}

	r.GET("/login", GetLoginForm)
	r.POST("/login", noteHandler.Login)

	r.GET("/register", GetRegistrationForm)
	r.POST("/register", noteHandler.CreateUser)

	r.GET("/", noteHandler.GetIndex)

	r.POST("/create", noteHandler.CreateNote)

	r.Run()
}
