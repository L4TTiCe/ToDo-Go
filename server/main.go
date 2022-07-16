package main

import (
	"log"
	"net/http"

	"github.com/L4TTiCe/ToDo-Go/server/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// loadEnv loads the .env file if it exists
func loadEnv() {
	log.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

// configureLogger sets optional flags for the logger
func configureLogger() {
	log.Print("Configuring logger...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initializeRouter() *gin.Engine {
	log.Println("Initializing router...")
	router := gin.Default()

	router.GET("/up", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is Up!")
	})

	return router
}

func main() {
	configureLogger()
	loadEnv()

	config.ConnectMongoDB()
	defer config.CloseClientDB()

	router := initializeRouter()

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
