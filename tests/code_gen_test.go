package tests

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/signup"
	"github.com/stretchr/testify/assert"
)

func TestGenerate4DigitCode(t *testing.T) {
	service := &signup.SignUpService{}
	code := service.GenerateCode(10)

	log := config.InitializeLogger().LogHandler
	log.Info(code)
	assert.Equal(t, 4, len(code))
}
