package controllers

// CRUD: user

import (
	"context"
	"net/http"

	"github.com/DroidZed/go_lance/config"
	"github.com/DroidZed/go_lance/db"
	"github.com/DroidZed/go_lance/models"
	"github.com/DroidZed/go_lance/utils"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"

)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	cursor, err := coll.Find(context.TODO(), nil)

	if err != nil {
		utils.JsonResponse(w, 400, utils.DtoResponse{Error: "Could not decode data!"})
	}

	var results []models.User = make([]models.User, 0)

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	utils.JsonResponse(w, 200, results)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	var result models.User = models.User{}

	err := coll.FindOne(context.TODO(), filter).Decode(result)

	if err != nil {
		utils.JsonResponse(w, 400, utils.DtoResponse{Error: "Could not decode data!"})
	}

	utils.JsonResponse(w, 200, result)

}
