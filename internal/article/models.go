package article

import (
	"github.com/DroidZed/my_blog/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Photo    string             `json:"photo" bson:"photo"`
	Tags     []string           `json:"tags" bson:"tags"`
	FileID   string             `json:"fileId" bson:"fileId"`
	AuthorID string             `json:"-" bson:"author_id"`
}

type ArticleWithUser struct {
	Article Article   `json:"article"`
	Author  user.User `json:"author"`
}
