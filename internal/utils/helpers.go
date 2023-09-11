package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
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

func GenerateAPICode() string {

	builder := strings.Builder{}

	const alpha = "A4BCD3EFG8HIJ6KLM7NO0PQRS2TUV9WX1YZ5"

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			builder.WriteByte(alpha[RNG(int64(j+17))])
		}
		if i != 3 {
			builder.WriteString("-")
		}
	}
	return strings.TrimSuffix(builder.String(), "-")
}
