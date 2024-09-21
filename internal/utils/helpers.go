package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type DtoResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type GetInput struct {
	Filter     bson.M
	Projection any
}

func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func GenUUID() string {

	code, err := uuid.NewRandom()

	if err != nil {
		panic(err)
	}

	return code.String()
}

/*
Decode the body to the appropriate format, infers type from usage.
*/
func DecodeBody[T interface{}](r *http.Request, out *T) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
