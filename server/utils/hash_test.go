package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHash(t *testing.T) {
	hash, err := BcryptHash("password")
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

func TestBcryptCheck(t *testing.T) {
	hash, err := BcryptHash("123456")
	assert.Nil(t, err)
	assert.True(t, BcryptCheck("123456", hash))
	assert.False(t, BcryptCheck("654321", hash))
}
