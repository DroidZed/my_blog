package tests

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/stretchr/testify/assert"
)

func TestGeneratingToken(t *testing.T) {

	token, err := cryptor.GenerateAccessToken("5s4dasd")

	if err != nil {
		t.Error(err)
	}

	t.Logf("=== DEBUG === token is: %s", token)

	assert.NotEmpty(t, token)
}
