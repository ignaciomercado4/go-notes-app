package main

import (
	"go-notes-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

func (h *NoteHandler) CreateUser(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBind(&authInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var user models.User
	h.DB.Where("username=?", authInput.Username).Find(&user)

	if user.ID != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already registered."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := models.User{
		Username: authInput.Username,
		Password: string(hashedPassword),
	}

	h.DB.Create(&newUser)

	c.JSON(http.StatusOK, gin.H{"data": newUser})
}

func (h *NoteHandler) Login(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBind(&authInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var user models.User

	h.DB.Where("username=?", authInput.Username).Find(&user)

	if user.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte("NOT-SECRET-KEY"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
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
