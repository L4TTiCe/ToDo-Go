package controller

import (
	"time"

	"github.com/L4TTiCe/ToDo-Go/server/models"
	"github.com/gin-gonic/gin"
)

func PopulateErrorResponse(c *gin.Context, errorResponse *models.ErrorResponse) {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	errorResponse.Timestamp = time.Now().UnixMilli()
	errorResponse.Path = scheme + "://" + c.Request.Host + c.Request.RequestURI
}
