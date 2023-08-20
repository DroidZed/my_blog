package controllers

import (
	"net/http"

	"github.com/DroidZed/go_lance/models"
	"github.com/DroidZed/go_lance/utils"
)

func GetDummyEntity(w http.ResponseWriter, r *http.Request) {

	dummyE := models.Dummy{Name: "Joe"}

	utils.JsonResponse(w, 200, dummyE)

}
