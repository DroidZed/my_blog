package tests

import (
	"github.com/DroidZed/go_lance/services"
	"testing"
)

const dummyPwd = "mysecretpassword"

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = services.HashPassword(dummyPwd)
	}
}

func BenchmarkCompareSecureToPlain(b *testing.B) {
	hashedPassword, _ := services.HashPassword(dummyPwd)

	for i := 0; i < b.N; i++ {
		_ = services.CompareSecureToPlain(hashedPassword, dummyPwd)
	}
}
