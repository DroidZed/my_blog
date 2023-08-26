package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const dummyPwd = "mysecrea56sd0a5s415asd415as5qaas]as'asda.da/*-as=-?/\\'sd;asda00..6as5d45-tpassword"

func TestCompareSecureToPlain(t *testing.T) {
	hashedPassword, err := HashPassword(dummyPwd)
	if err != nil {
		t.Fatal(err)
	}

	result := CompareSecureToPlain(hashedPassword, dummyPwd)

	assert.Equal(t, true, result)
}

func BenchmarkHashPassword(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		HashPassword(dummyPwd)
	}
}
