package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDefault(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(Default())
}
