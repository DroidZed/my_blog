package cryptor_test

import (
	"testing"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/stretchr/testify/assert"
)

const dummyPwd = "mysecrea56sd0a5s415asd415as5qaas]as'asda.da/*-as=-?/\\'sd;"

func TestCompareSecureToPlain(t *testing.T) {

	c := &cryptor.Cryptor{}

	hashedPassword, err := c.HashPlain(dummyPwd)
	if err != nil {
		t.Fatal(err)
	}

	result := c.CompareSecureToPlain(hashedPassword, dummyPwd)

	if result != true {
		t.Error("failed comparison, passwords don't match")
	}
}

func BenchmarkHashPassword(b *testing.B) {
	c := &cryptor.Cryptor{}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.HashPlain(dummyPwd)
	}
}

func TestGeneratingToken(t *testing.T) {

	tokenor := &cryptor.Cryptor{}

	token, err := tokenor.GenerateAccessToken("64f82060217cdc32997cc7b3")

	if err != nil {
		t.Error(err)
	}

	assert.NotEmpty(t, token)

	if token == "" {
		t.Error("GenerateAccessToken() got empty, expected value")
	}
}

func BenchmarkGeneratingTokens(b *testing.B) {
	tokenor := &cryptor.Cryptor{}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tokenor.GenerateAccessToken("64f82060217cdc32997cc7b3")
	}
}
