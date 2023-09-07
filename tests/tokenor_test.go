package tests

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/stretchr/testify/assert"
)

func TestGeneratingToken(t *testing.T) {

	token, err := cryptor.GenerateAccessToken("64f82060217cdc32997cc7b3")

	if err != nil {
		t.Error(err)
	}

	t.Logf("=== DEBUG === token is: %s", token)

	assert.NotEmpty(t, token)
}

func BenchmarkGeneratingTokens(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cryptor.GenerateAccessToken("64f82060217cdc32997cc7b3")
	}
}
