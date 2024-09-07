package forgotPwd

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForgotPwdReq struct {
	Email string `json:"email"`
}

type MagicCode struct {
	ID    primitive.ObjectID `bson:"_id" json:"-"`
	Email string             `json:"email" bson:"email"`
	Token string             `json:"token" bson:"token"`
}
