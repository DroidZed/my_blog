package user

import (
	"github.com/DroidZed/my_blog/internal/cryptor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"`
	FullName string             `json:"fullName" bson:"fullName"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password,omitempty"`
	Photo    string             `json:"photo" bson:"photo,omitempty"`
}

func (u *User) HashUserPassword() (*User, error) {

	hashed, err := cryptor.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashed

	return u, nil
}
