package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/up", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is Up!")
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
