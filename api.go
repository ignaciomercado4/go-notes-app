package main

import (
	"go-notes-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteHandler struct {
	DB *gorm.DB
}

// User auth
func GetRegistrationForm(c *gin.Context) {
	c.HTML(http.StatusOK, "registrationForm.tmpl", gin.H{
		"title": "Register",
	})
}

func GetLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "loginForm.tmpl", gin.H{
		"title": "login",
	})
}

// Index
func (h *NoteHandler) GetIndex(c *gin.Context) {
	var notes []models.Note

	result := h.DB.Find(&notes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"notes": notes,
	})
}

// Notes
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var newNote models.Note

	if err := c.ShouldBind(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Binding error",
			"details": err.Error(),
		})

		return
	}

	result := h.DB.Create(&newNote)

	if result != nil {
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	c.Redirect(http.StatusSeeOther, "/")
}
