package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIndex(c *gin.Context) {

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Notitas",
	})
}
