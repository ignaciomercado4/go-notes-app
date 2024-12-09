package main

import (
	"go-notes-app/database"
	middlewares "go-notes-app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.ConnectDatabase()

	r.LoadHTMLGlob("templates/*")

	r.Use(middlewares.CheckAuth(db))

	noteHandler := &NoteHandler{DB: db}

	r.GET("/login", GetLoginForm)
	r.POST("/login", noteHandler.Login)

	r.GET("/register", GetRegistrationForm)
	r.POST("/register", noteHandler.CreateUser)

	r.GET("/logout", noteHandler.Logout)

	r.GET("/", noteHandler.GetIndex)

	r.POST("/create", noteHandler.CreateNote)
	r.GET("/delete/:id", noteHandler.DeleteNote)

	r.Run()
}
