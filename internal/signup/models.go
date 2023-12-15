package signup

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConfirmationCode struct {
	ID        primitive.ObjectID `bson:"_id"`
	Code      string             `bson:"code"`
	ExpiresAt primitive.DateTime `bson:"expiresAt"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	Email     string             `bson:"email"`
}
