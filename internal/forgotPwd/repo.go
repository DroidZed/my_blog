package forgotPwd

import "go.mongodb.org/mongo-driver/mongo"

type MagicTokenProvider interface{}

type Repo struct {
	DbClient *mongo.Client
	DBName   string
}
