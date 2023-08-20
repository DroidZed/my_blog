package controllers

// CRUD: user

import (
	"context"
	"fmt"
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
		return
	}

	results := make([]models.User, 0)

	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.JsonResponse(w, 400, utils.DtoResponse{Error: fmt.Sprintf("Could not decode data!\n%s", err.Error())})
		return
	}

	defer cursor.Close(context.TODO())

	utils.JsonResponse(w, 200, results)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	var result models.User = models.User{}

	err := coll.FindOne(context.TODO(), filter).Decode(result)

	if err != nil {
		utils.JsonResponse(w, 400, utils.DtoResponse{Error: fmt.Sprintf("Could not decode data!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, 200, result)

}
