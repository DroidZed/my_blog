package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	FullName string             `json:"fullName"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	photo    string             `json:"photo"`
	age      int                `json:"age"`
}
