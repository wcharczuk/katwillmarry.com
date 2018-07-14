package util

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAnyError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(AnyError())
	assert.Nil(AnyError(nil))
	assert.Nil(AnyError(nil, nil))
	assert.NotNil(AnyError(fmt.Errorf("test")))
	assert.NotNil(AnyError(nil, fmt.Errorf("test")))
	assert.NotNil(AnyError(nil, fmt.Errorf("test"), nil))
}
