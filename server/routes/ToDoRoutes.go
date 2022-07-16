package routes

import (
	"github.com/L4TTiCe/ToDo-Go/server/controller/ToDoItemController"
	"github.com/gin-gonic/gin"
)

func ToDoRoute(router *gin.Engine) {
	routerGroup := router.Group("/todo")

	routerGroup.GET("/up", ToDoItemController.HealthCheck)
}
