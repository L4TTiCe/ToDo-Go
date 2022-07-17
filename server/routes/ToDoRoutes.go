package routes

import (
	"github.com/L4TTiCe/ToDo-Go/server/controller/ToDoItemController"
	"github.com/gin-gonic/gin"
)

// ToDoRoutes contains the routes for the ToDo API.
func ToDoRoutes(router *gin.Engine) {
	routerGroup := router.Group("/todo")

	routerGroup.GET("/up", ToDoItemController.HealthCheck)

	routerGroup.POST("/", ToDoItemController.Create)
	routerGroup.GET("/", ToDoItemController.RetrieveAll)
	routerGroup.GET("/:id", ToDoItemController.RetrieveOne)
	routerGroup.PUT("/:id", ToDoItemController.UpdateOne)
	routerGroup.DELETE("/:id", ToDoItemController.DeleteOne)
}
