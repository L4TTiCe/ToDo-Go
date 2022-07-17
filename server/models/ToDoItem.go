package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ToDoItem is a struct that contains the ToDoItem data.
// CreatedAt is a timestamp that are automatically set when the ToDoItem is created, and is represented as a Unix millisecond timestamp.
// Similarly, deadline is an optional timestamp that represents the deadline of the ToDoItem.
type ToDoItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed,omitempty"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt,omitempty"`
	Deadline  int64              `bson:"deadline,omitempty" json:"deadline,omitempty"`
}
