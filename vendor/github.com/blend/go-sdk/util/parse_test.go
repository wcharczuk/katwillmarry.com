package util

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	assert.Zero(Parse.Float64("foo"))
	assert.Equal(3.14, Parse.Float64("3.14"))

	assert.Zero(Parse.Float32("foo"))
	assert.Equal(3.14, Parse.Float32("3.14"))

	assert.Zero(Parse.Int("foo"))
	assert.Equal(3, Parse.Int("3"))

	assert.Zero(Parse.Int32("foo"))
	assert.Equal(3, Parse.Int32("3"))

	assert.Zero(Parse.Int64("foo"))
	assert.Equal(3, Parse.Int64("3"))

	assert.Empty(Parse.Ints())
	values, err := Parse.Ints("1", "2", "3")
	assert.Nil(err)
	assert.Equal([]int{1, 2, 3}, values)

	assert.Empty(Parse.Int64s())
	int64values, err := Parse.Int64s("1", "2", "3")
	assert.Nil(err)
	assert.Equal([]int64{1, 2, 3}, int64values)

	assert.Empty(Parse.Float64s())
	float64values, err := Parse.Float64s("1.1", "2.2", "3.3")
	assert.Nil(err)
	assert.Equal([]float64{1.1, 2.2, 3.3}, float64values)
}
