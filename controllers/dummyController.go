package controllers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"
//

// 	"github.com/DroidZed/go_lance/models"
// 	"github.com/DroidZed/go_lance/utils"
// 	"github.com/go-chi/chi"

// )

// var dummies = make([]models.Dummy, 5)

// func GetDummyEntity(w http.ResponseWriter, r *http.Request) {

// 	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

// 	if err != nil {
// 		utils.JsonResponse(w, 404, utils.DtoResponse{Error: "Not Found"})
// 	}

// 	var found *models.Dummy = nil

// 	for i := 0; i < len(dummies); i++ {
// 		if dummies[i].Id == int(id) {
// 			*found = dummies[i]
// 			break
// 		}
// 	}

// 	utils.JsonResponse(w, http.StatusOK, found)

// }

// func AddDummy(w http.ResponseWriter, r *http.Request) {

// 	d := &models.Dummy{}

// 	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
// 		log.Fatal(err)
// 	}

// 	dummies = append(dummies, *d)

// 	utils.JsonResponse(w, 200, dummies)

// }

// func UpdateDummy(w http.ResponseWriter, r *http.Request) {

// 	d := &models.Dummy{}

// 	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
// 		log.Fatal(err)
// 	}

// 	index := 0

// 	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

// 	if err != nil {
// 		utils.JsonResponse(w, 404, utils.DtoResponse{Error: "Not Found"})
// 	}

// 	for i:= 0; i < len(dummies); i++ {
// 		if dummies[i].Id == int(id) {
// 			index = i
// 		}
// 	}

// 	dummies[index].Name = d.Name

// 	utils.JsonResponse(w, 200, dummies)

// }
