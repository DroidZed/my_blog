package tests

import (
	"testing"
	"time"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/stretchr/testify/assert"
)

func TestGeneratingToken(t *testing.T) {

	claims := make(map[string]any)

	claims["iat"] = time.Now().Unix()
	claims["sub"] = "5s4dasd"

	token, err := cryptor.GenerateAccessToken(claims)

	if err != nil {
		t.Error(err)
	}

	t.Logf("=== DEBUG === token is: %s", token)

	assert.NotEmpty(t, token)
}
