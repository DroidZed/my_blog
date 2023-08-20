package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("All users here"))

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	w.Write([]byte(fmt.Sprintf("User with id: %s", id)))
}
