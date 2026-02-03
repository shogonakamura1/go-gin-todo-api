package main

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Status": http.StatusOK,
		})
	})
	r.Run()
}
