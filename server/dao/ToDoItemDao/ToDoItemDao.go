package ToDoItemDao

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/L4TTiCe/ToDo-Go/server/config"
	"github.com/L4TTiCe/ToDo-Go/server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create creates a new ToDoItem in the DB.
// It takes a ToDoItem struct and returns a struct with the InsertedID or an ErrorResponse.
func Create(item *models.ToDoItem) (interface{}, *models.ErrorResponse) {
	// Check if item is nil
	if item == nil {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Empty Request",
			Detail: "Request body is empty",
		}
	}

	// Check if Title is empty
	if item.Title == "" {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "No Title",
			Detail: "Title is a required field",
		}
	}

	// Overwrite CreatedAt field with current server time
	item.CreatedAt = time.Now().UnixMilli()

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert item into DB
	result, err := config.ToDoItemsCollection.InsertOne(ctx, &item)
	if err != nil {
		return nil, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
		}
	}

	return result, nil
}

func RetriveAll(sortParam string, sortOrder int) ([]models.ToDoItem, *models.ErrorResponse) {
	log.Println("ToDo: RetriveAll (sortParam: " + sortParam + ", sortOrder: " + strconv.Itoa(sortOrder) + ")")

	// Validate sortParam
	if sortParam != "" && sortParam != "title" && sortParam != "completed" && sortParam != "createdAt" && sortParam != "deadline" {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Invalid Sort Parameter",
			Detail: "Sort parameter must be one of the following: title, completed, createdAt, deadline",
		}
	}

	// Validate sortOrder
	if sortOrder != 1 && sortOrder != -1 {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Invalid Sort Order",
			Detail: "Sort parameter must be one of the following: asc, desc, 1, -1",
		}
	}

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create options for sorting
	// Note: bson.D{} preserves order and is ideal for specifing sort ordering
	cursor, err := config.ToDoItemsCollection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: sortParam, Value: sortOrder}}))
	if err != nil {
		log.Print(err)
		return nil, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
		}
	}

	// Create a slice to hold all ToDoItems
	var items []models.ToDoItem

	// Append items from cursor to items
	for cursor.Next(context.Background()) {
		item := models.ToDoItem{}
		err = cursor.Decode(&item)
		if err != nil {
			log.Print(err)
			return nil, &models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  err.Error(),
			}
		}
		items = append(items, item)
	}

	return items, nil
}
