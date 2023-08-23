package controllers

// CRUD: user

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/DroidZed/go_lance/config"
	"github.com/DroidZed/go_lance/db"
	"github.com/DroidZed/go_lance/models"
	"github.com/DroidZed/go_lance/services"
	"github.com/DroidZed/go_lance/utils"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	log := services.Logger.LogHandler

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := coll.Find(ctx, bson.D{})

	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: "Could not decode data!"})
		return
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Errorf("Error closing cursor.\n%s", err.Error())
		}
	}(cur, ctx)

	results := make([]models.User, 0)

	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Could not decode data!\n%s", err.Error())})
			return
		}
	}

	if err := cur.Err(); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, utils.DtoResponse{Error: fmt.Sprintf("Runtime error:\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, results)

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := models.User{}

	err := coll.FindOne(ctx, filter).Decode(result)

	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Could not decode data!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, result)

}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	coll := db.Client.Database(config.EnvDbName()).Collection("users")

	filter := bson.D{{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := coll.DeleteOne(ctx, filter)

	if err != nil || result.DeletedCount == 0 {
		utils.JsonResponse(w, http.StatusBadRequest, utils.DtoResponse{Error: fmt.Sprintf("Could not delete user!\n%s", err.Error())})
		return
	}

	utils.JsonResponse(w, http.StatusOK, utils.DtoResponse{Message: fmt.Sprintf("User with id: %s has been deleted successfully!", id)})
}
