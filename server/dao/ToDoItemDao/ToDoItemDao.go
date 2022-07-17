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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func RetrieveAll(sortParam string, sortOrder int) ([]models.ToDoItem, *models.ErrorResponse) {
	log.Println("ToDo: RetrieveAll (sortParam: " + sortParam + ", sortOrder: " + strconv.Itoa(sortOrder) + ")")

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

func RetrieveAllWithParams(attrib string, verb string, date int64, sortOrder int) ([]models.ToDoItem, *models.ErrorResponse) {
	log.Println("ToDo: RetrieveAllWithParams (attrib: " + attrib + ", verb: " + verb + ", date: " + strconv.FormatInt(date, 10) + ")")

	// Validate attrib
	if attrib != "createdAt" && attrib != "deadline" {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Invalid Attribute",
			Detail: "Attribute must be one of the following: createdAt, deadline",
		}
	}

	// Validate verb
	if verb != "gte" && verb != "lte" {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Invalid Verb",
			Detail: "Verb must be one of the following: gte, lte",
		}
	}

	// Validate date
	if date < 0 {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  "Invalid Date",
			Detail: "Date must be a positive integer",
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
	sortOptions := options.Find().SetSort(bson.D{{Key: attrib, Value: sortOrder}})

	// Create a filter for the query
	filter := bson.D{{Key: attrib, Value: bson.D{{Key: "$" + verb, Value: date}}}}

	cursor, err := config.ToDoItemsCollection.Find(ctx, filter, sortOptions)
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

// RetriveOne retrieves a ToDoItem from the DB.
// It takes an ID and returns a ToDoItem or an ErrorResponse.
func RetriveOne(id string) (*models.ToDoItem, *models.ErrorResponse) {
	log.Print("ToDo: RetriveOne (id: " + id + ")")

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  err.Error(),
			Detail: "Invalid ID",
		}
	}

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find item in DB
	item := models.ToDoItem{}
	err = config.ToDoItemsCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&item)
	if err != nil {
		// Check if item was not found
		if err.Error() == "mongo: no documents in result" {
			return nil, &models.ErrorResponse{
				Status: http.StatusNotFound,
				Title:  "Item Not Found",
				Detail: "Item with ID " + id + " not found",
			}
		} else {
			log.Print(err)
			return nil, &models.ErrorResponse{
				Status: http.StatusInternalServerError,
				Title:  err.Error(),
			}
		}
	}

	return &item, nil
}

// Update updates a ToDoItem in the DB.
// It takes a ToDoItem struct and returns the update status or an ErrorResponse.
func UpdateOne(id string, updatedItem *models.ToDoItem) (interface{}, *models.ErrorResponse) {
	log.Print("ToDo: UpdateOne (id: " + id + ")")

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  err.Error(),
			Detail: "Invalid ID",
		}
	}

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Item validation already performed when FindOne is called from contrller before this function is called

	// Update item in DB
	result, err := config.ToDoItemsCollection.ReplaceOne(ctx, bson.M{"_id": objectId}, &updatedItem)
	if err != nil {
		log.Print(err)
		return nil, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
		}
	}

	return result, nil
}

// DeleteOne deletes a ToDoItem from the DB.
// It takes an ID and returns the status of the delete Operation or an ErrorResponse.
func DeleteOne(id string) (interface{}, *models.ErrorResponse) {
	log.Print("ToDo: DeleteOne (id: " + id + ")")

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Title:  err.Error(),
			Detail: "Invalid ID",
		}
	}

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete item in DB
	result, err := config.ToDoItemsCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		log.Print(err)
		return nil, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Title:  err.Error(),
		}
	}

	// Check if item was deleted
	if result.DeletedCount == 0 {
		return nil, &models.ErrorResponse{
			Status: http.StatusNotFound,
			Title:  "Item Not Found",
			Detail: "Item with ID " + id + " not found",
		}
	}

	return result, nil
}
