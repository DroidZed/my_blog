package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"unsafe"
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

func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func RNG(outerBound int64) int64 {
	if outerBound <= 0 {
		return -1
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(outerBound))
	if err != nil {
		return -1
	}

	n := nBig.Int64()

	if n == 0 {
		return n + 1
	}
	return n
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

func DecodeBody[T interface{}](r *http.Request, out *T) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
