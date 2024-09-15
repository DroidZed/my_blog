package utils_test

import (
	"testing"

	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenApiKey(t *testing.T) {

	code := utils.GenUUID()

	assert.NotEmpty(t, code)

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
