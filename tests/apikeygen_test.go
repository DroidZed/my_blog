package tests

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenApiKey(t *testing.T) {

	code := utils.GenerateAPICode()

	log := config.InitializeLogger().LogHandler

	log.Debugf("Code: %s", code)

	assert.NotEmpty(t, code)
	last := code[len(code)-1]
	assert.NotEqual(t, "-", last)
}

func BenchmarkGenApiKey(b *testing.B) {

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		utils.GenerateAPICode()
	}
}
