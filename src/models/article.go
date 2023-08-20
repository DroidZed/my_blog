package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Article struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Title   string             `json:"title"`
	content string             `json:"content"`
	Link    string             `json:"link"`
	photo   string             `json:"photo"`
	Tags    []string           `json:"tags"`
}
