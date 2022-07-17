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
	start := c.Request.URL.Query().Get("start")
	end := c.Request.URL.Query().Get("end")

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

		var startDate int64
		var endDate int64

		// Either before or after must be specified, together with a verb
		if before != "" && after == "" && start == "" && end == "" { // before
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
		} else if after != "" && before == "" && start == "" && end == "" { // after
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
		} else if after != "" && before != "" && start == "" && end == "" { // after and before
			errorResponse = &models.ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Invalid Query",
				Detail: "Must specify either before or after",
			}
			controller.PopulateErrorResponse(c, errorResponse)
			c.JSON(errorResponse.Status, errorResponse)
			return
		} else if (after != "" || before != "") && start != "" && end != "" { // start and end and after or before
			errorResponse = &models.ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Invalid Query",
				Detail: "Must specify either before / after or start and end",
			}
			controller.PopulateErrorResponse(c, errorResponse)
			c.JSON(errorResponse.Status, errorResponse)
			return
		} else if after == "" && before == "" && start != "" && end != "" { // start and end
			// validate start and end
			intVal, err := strconv.Atoi(start)
			if err != nil {
				errorResponse = &models.ErrorResponse{
					Status: http.StatusBadRequest,
					Title:  "Invalid Date",
					Detail: "Date must be a positive integer. Check start.",
				}
				controller.PopulateErrorResponse(c, errorResponse)
				c.JSON(errorResponse.Status, errorResponse)
				return
			}

			startDate = int64(intVal)

			intVal, err = strconv.Atoi(end)
			if err != nil {
				errorResponse = &models.ErrorResponse{
					Status: http.StatusBadRequest,
					Title:  "Invalid Date",
					Detail: "Date must be a positive integer. Check end.",
				}
				controller.PopulateErrorResponse(c, errorResponse)
				c.JSON(errorResponse.Status, errorResponse)
				return
			}

			endDate = int64(intVal)
		}

		if before == "" && after == "" && start == "" && end == "" { // no params (before, after, start, end)
			result, errorResponse = ToDoItemDao.RetrieveAll(attrib, sortOrder)
		} else if (before != "" || after != "") && start == "" && end == "" { // before or after
			result, errorResponse = ToDoItemDao.RetrieveWithParams(attrib, verb, date, sortOrder)
		} else if start != "" && end != "" { // start and end
			result, errorResponse = ToDoItemDao.RetrieveBetween(attrib, startDate, endDate, sortOrder)
		}

	} else {
		if before != "" || after != "" || start != "" || end != "" {
			errorResponse = &models.ErrorResponse{
				Status: http.StatusBadRequest,
				Title:  "Invalid Query",
				Detail: "Must specify an attribute to use with the before, after, start, or end parameters",
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

	result, errorResponse = ToDoItemDao.RetrieveOne(id)

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

	item, errorResponse := ToDoItemDao.RetrieveOne(id)

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

func DeleteOne(c *gin.Context) {
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
