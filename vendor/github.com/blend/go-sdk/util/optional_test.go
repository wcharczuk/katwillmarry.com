package util

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptional(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(OptionalUInt8(1))
	assert.NotNil(OptionalUInt16(1))
	assert.NotNil(OptionalUInt(1))
	assert.NotNil(OptionalUInt64(1))
	assert.NotNil(OptionalInt16(1))
	assert.NotNil(OptionalInt(1))
	assert.NotNil(OptionalInt32(1))
	assert.NotNil(OptionalInt64(1))
	assert.NotNil(OptionalFloat32(1))
	assert.NotNil(OptionalFloat64(1))
	assert.NotNil(OptionalString("1"))
	assert.NotNil(OptionalBool(true))
	assert.NotNil(OptionalTime(time.Time{}))
	assert.NotNil(OptionalDuration(time.Second))
}
