package utils_test

import (
	"testing"

	"github.com/DroidZed/my_blog/internal/utils"
)

func TestGenApiKey(t *testing.T) {

	code := utils.GenUUID()

	if code == "" {
		t.Errorf("GenerateAPICode() got %v, want not empty", code)
	}

	last := code[len(code)-1]

	if last == '-' {
		t.Errorf("GenerateAPICode() got -, want not %v", last)
	}
}

func BenchmarkGenApiKey(b *testing.B) {

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		utils.GenUUID()
	}
}
