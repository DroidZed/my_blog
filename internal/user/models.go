package user

import (
	"github.com/DroidZed/go_lance/internal/cryptor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConfirmationToken struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Token  string             `bson:"token" json:"token"`
	UserId string             `bson:"userId" json:"userId"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	FullName  string             `json:"fullName" bson:"fullName"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Photo     string             `json:"photo" bson:"photo"`
	Age       int                `json:"age" bson:"age"`
	Role      string             `json:"role,omitempty" bson:"role"`
	AccStatus int                `json:"accStatus,omitempty" bson:"accStatus"`
}

func (u *User) HashUserPassword() (*User, error) {

	hashed, err := cryptor.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashed

	return u, nil
}
