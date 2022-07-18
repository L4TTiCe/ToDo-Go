package main

import (
	"log"
	"net/http"
	"os"

	"github.com/L4TTiCe/ToDo-Go/server/config"
	"github.com/L4TTiCe/ToDo-Go/server/routes"
	"github.com/gin-contrib/cors"

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
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initializeRouter() *gin.Engine {
	log.Println("Initializing router...")
	router := gin.Default()

	// Enable CORS for all requests
	router.Use(cors.Default())

	healthCheck := func(c *gin.Context) {
		c.String(http.StatusOK, "Server is Up!")
	}

	router.GET("/up", healthCheck)

	routes.ToDoRoutes(router)

	return router
}

func main() {
	configureLogger()
	loadEnv()

	config.ConnectMongoDB()
	defer config.CloseClientDB()

	router := initializeRouter()

	err := router.Run()
	if err != nil {
		panic(err)
	}
}
