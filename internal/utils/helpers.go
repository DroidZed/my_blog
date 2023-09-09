package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
)

type DtoResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
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

func StringToBytes(s string) []byte {

	bytes := make([]byte, len(s))
	defer func() {
		bytes = nil
	}()

	copy(bytes, s)

	return bytes
}

func RNG(outerBound int64) int64 {
	if outerBound <= 0 {
		return -1
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(outerBound))
	if err != nil {
		return -1
	}

	return nBig.Int64()
}
