package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToDoItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Completed int64              `bson:"completed" json:"completed,omitempty"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt,omitempty"`
	Deadline  int64              `bson:"deadline,omitempty" json:"deadline,omitempty"`
}
