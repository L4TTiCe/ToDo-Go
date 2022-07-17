package ToDoItemController

import (
	"net/http"
	"strconv"

	"github.com/L4TTiCe/ToDo-Go/server/controller"
	"github.com/L4TTiCe/ToDo-Go/server/dao/ToDoItemDao"
	"github.com/L4TTiCe/ToDo-Go/server/models"
	"github.com/gin-gonic/gin"
)

// HealthCheck is a handler function that returns a 200 response.
func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ToDoRoute is Up!")
}

// Create is a handler function that creates a new ToDoItem.
// It takes a JSON body and creates a new ToDoItem in the DB.
// It returns a JSON response with the newly created ToDoItem's ID or an error.
func Create(c *gin.Context) {
	var item models.ToDoItem

	// Bind JSON to struct
	err := c.BindJSON(&item)
	if err != nil {
		errorResponse := &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
			Detail: "Error parsing JSON",
		}
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	// Attempt to create item in DB using DAO
	result, errorResponse := ToDoItemDao.Create(&item)
	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	// Populate success response before sending to client
	c.JSON(http.StatusCreated, &result)
}

func RetrieveAll(c *gin.Context) {
	attrib := c.Request.URL.Query().Get("attrib")
	sort := c.Request.URL.Query().Get("sort")
	before := c.Request.URL.Query().Get("before")
	after := c.Request.URL.Query().Get("after")

	var sortOrder int

	switch sort {
	case "asc", "1":
		sortOrder = 1
	case "desc", "-1":
		sortOrder = -1
	case "":
		sortOrder = 1
	default:
		sortOrder = 0
	}

	var result []models.ToDoItem
	var errorResponse *models.ErrorResponse

	if attrib != "" {
		var verb string
		var date int64

		if before != "" && after == "" {
			verb = "lte"
			intVal, err := strconv.Atoi(before)
			if err != nil {
				errorResponse = &models.ErrorResponse{
					Status: http.StatusBadRequest,
					Title:  "Invalid Date",
					Detail: "Date must be a positive integer",
				}
				controller.PopulateErrorResponse(c, errorResponse)
				c.JSON(errorResponse.Status, errorResponse)
				return
			}

			date = int64(intVal)
		} else if after != "" && before == "" {
			verb = "gte"
			intVal, err := strconv.Atoi(after)
			if err != nil {
				errorResponse = &models.ErrorResponse{
					Status: http.StatusBadRequest,
					Title:  "Invalid Date",
					Detail: "Date must be a positive integer",
				}
				controller.PopulateErrorResponse(c, errorResponse)
				c.JSON(errorResponse.Status, errorResponse)
				return
			}

			date = int64(intVal)
		} else if after != "" && before != "" {
			errorResponse = &models.ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Invalid Query",
				Detail: "Must specify either before or after",
			}
			controller.PopulateErrorResponse(c, errorResponse)
			c.JSON(errorResponse.Status, errorResponse)
			return
		}

		if before == "" && after == "" {
			result, errorResponse = ToDoItemDao.RetrieveAll(attrib, sortOrder)
		} else {
			result, errorResponse = ToDoItemDao.RetrieveAllWithParams(attrib, verb, date, sortOrder)
		}

	} else {
		if before != "" || after != "" {
			errorResponse = &models.ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Invalid Query",
				Detail: "Must specify attrib to use before or after",
			}
			controller.PopulateErrorResponse(c, errorResponse)
			c.JSON(errorResponse.Status, errorResponse)
			return
		}

		result, errorResponse = ToDoItemDao.RetrieveAll("createdAt", sortOrder)
	}

	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	c.JSON(http.StatusOK, result)
}

func RetrieveOne(c *gin.Context) {
	id := c.Param("id")

	var result *models.ToDoItem
	var errorResponse *models.ErrorResponse

	result, errorResponse = ToDoItemDao.RetriveOne(id)

	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	c.JSON(http.StatusOK, &result)
}

func UpdateOne(c *gin.Context) {
	id := c.Param("id")

	item, errorResponse := ToDoItemDao.RetriveOne(id)

	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	// Bind JSON to struct
	err := c.BindJSON(&item)
	if err != nil {
		errorResponse = &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
			Detail: "Error parsing JSON",
		}
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	result, errorResponse := ToDoItemDao.UpdateOne(id, item)

	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	c.JSON(http.StatusOK, &result)
}

func DeleteeOne(c *gin.Context) {
	id := c.Param("id")

	var errorResponse *models.ErrorResponse

	result, errorResponse := ToDoItemDao.DeleteOne(id)

	if errorResponse != nil {
		// Populate error response before sending to client
		controller.PopulateErrorResponse(c, errorResponse)

		c.JSON(errorResponse.Status, errorResponse)
		return
	}

	c.JSON(http.StatusOK, &result)
}
