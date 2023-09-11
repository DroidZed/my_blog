package register

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConfirmationToken struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Token  string             `bson:"token" json:"token"`
	UserId string             `bson:"userId" json:"userId"`
}
