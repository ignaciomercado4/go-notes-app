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
		c.HTML(http.StatusBadRequest, "registrationForm.tmpl", gin.H{
			"title": "Register",
			"error": "Error en el formulario",
		})
		return
	}

	var user models.User
	h.DB.Where("username=?", authInput.Username).Find(&user)

	if user.ID != 0 {
		c.HTML(http.StatusBadRequest, "registrationForm.tmpl", gin.H{
			"title": "Register",
			"error": "Username already registered.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusBadRequest, "registrationForm.tmpl", gin.H{
			"title": "Register",
			"error": "Error hashing password",
		})
		return
	}

	newUser := models.User{
		Username: authInput.Username,
		Password: string(hashedPassword),
	}

	h.DB.Create(&newUser)

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  newUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte("NOT-SECRET-KEY"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "registrationForm.tmpl", gin.H{
			"title": "Register",
			"error": "Error generating token",
		})
		return
	}

	c.SetCookie("Authorization", "Bearer "+token, 3600, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *NoteHandler) Login(c *gin.Context) {
	var authInput models.AuthInput

	if err := c.ShouldBind(&authInput); err != nil {
		c.HTML(http.StatusBadRequest, "loginForm.tmpl", gin.H{
			"title": "Login",
			"error": "Error en el formulario",
		})
		return
	}

	var user models.User
	h.DB.Where("username=?", authInput.Username).Find(&user)

	if user.ID == 0 {
		c.HTML(http.StatusUnauthorized, "loginForm.tmpl", gin.H{
			"title": "Login",
			"error": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "loginForm.tmpl", gin.H{
			"title": "Login",
			"error": "Invalid password",
		})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte("NOT-SECRET-KEY"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "loginForm.tmpl", gin.H{
			"title": "Login",
			"error": "Error generating token",
		})
		return
	}

	c.SetCookie("Authorization", "Bearer "+token, 3600, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/")
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
