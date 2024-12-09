package middleware

import (
	"fmt"
	"go-notes-app/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func CheckAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero intenta obtener el token de la cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			// Si no hay cookie, intenta obtener el token del header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				// Rutas públicas que no requieren autenticación
				publicRoutes := []string{"/login", "/register", "/css/", "/js/"}

				currentPath := c.Request.URL.Path
				isPublicRoute := false
				for _, route := range publicRoutes {
					if strings.HasPrefix(currentPath, route) {
						isPublicRoute = true
						break
					}
				}

				if !isPublicRoute {
					c.Redirect(http.StatusTemporaryRedirect, "/login")
					c.Abort()
					return
				}
				return
			}
			tokenString = authHeader
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		var user models.User
		db.Where("ID=?", claims["id"]).Find(&user)

		if user.ID == 0 {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
