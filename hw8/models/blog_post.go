package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogPost model
type BlogPost struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string
	Date    string
	Link    string
	Content string
}
