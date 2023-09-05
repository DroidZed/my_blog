package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/go-chi/chi/v5"
)

type DtoResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SetupHostWithPort(host string, port int64) string { return fmt.Sprintf("%s:%d", host, port) }

func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func LogAllRoutes(r chi.Routes) {

	log := config.InitializeLogger().LogHandler

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Debugf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Errorf("Logging err: %s\n", err.Error())
	}
}

func StringToBytes(s string) []byte {

	bytes := make([]byte, len(s)+10)
	defer func() {
		bytes = nil
	}()

	copy(bytes, s)

	return bytes
}

func LinearRandomGenerator(seed int64, multiplier int64, increment int64, modulus int64, length int) []int64 {

	results := make([]int64, 0)
	for i := 0; i < length; i++ {
		seed = (multiplier*seed + increment) % modulus
		results = append(results, seed)
	}

	return results
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
