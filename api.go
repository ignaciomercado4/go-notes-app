package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteHandler struct {
	DB *gorm.DB
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var newNote Note

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

func (h *NoteHandler) GetIndex(c *gin.Context) {
	var notes []Note

	result := h.DB.Find(&notes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"notes": notes,
	})
}
