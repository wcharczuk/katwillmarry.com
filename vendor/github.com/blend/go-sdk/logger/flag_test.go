package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAsFlags(t *testing.T) {
	assert := assert.New(t)

	assert.Len(AsStrings(Error, Warning, Fatal), 3)
	assert.Len(AsFlags("Error", "Warning", "Fatal"), 3)

	flags := []string{"foo", "bar", "baz", "buzz"}
	assert.Equal(flags, AsStrings(AsFlags(flags...)...))
}
