package article

import "go.mongodb.org/mongo-driver/bson/primitive"

type Article struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Link     string             `json:"link" bson:"link"`
	Photo    string             `json:"photo" bson:"photo"`
	Tags     []string           `json:"tags" bson:"tags"`
	FileID   string             `json:"fileId" bson:"fileId"`
	AuthorId string             `json:"authorId" bson:"author_id"`
}
