package util

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCoalesceString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", Coalesce.String("", ""))
	assert.Equal("foo", Coalesce.String("", "foo"))
	assert.Equal("bar", Coalesce.String("", "foo", "bar"))
	assert.Equal("bar", Coalesce.String("", "foo", "bar", "baz"))
	assert.Equal("moo", Coalesce.String("moo", "foo", "bar", "baz"))
}

func TestCoalesceBool(t *testing.T) {
	assert := assert.New(t)

	assert.False(Coalesce.Bool(nil, false))
	assert.True(Coalesce.Bool(nil, true))
	assert.False(Coalesce.Bool(nil, true, false))
	assert.True(Coalesce.Bool(OptionalBool(true), false, false))
}

func TestCoalesceInt(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Int(0, 1))
	assert.Equal(2, Coalesce.Int(0, 1, 2))
	assert.Equal(2, Coalesce.Int(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Int(4, 1, 2, 3))
}

func TestCoalesceInt32(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Int32(0, 1))
	assert.Equal(2, Coalesce.Int32(0, 1, 2))
	assert.Equal(2, Coalesce.Int32(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Int32(4, 1, 2, 3))
}

func TestCoalesceInt64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Int64(0, 1))
	assert.Equal(2, Coalesce.Int64(0, 1, 2))
	assert.Equal(2, Coalesce.Int64(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Int64(4, 1, 2, 3))
}

func TestCoalesceFloat32(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Float32(0, 1))
	assert.Equal(2, Coalesce.Float32(0, 1, 2))
	assert.Equal(2, Coalesce.Float32(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Float32(4, 1, 2, 3))
}

func TestCoalesceFloat64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Float64(0, 1))
	assert.Equal(2, Coalesce.Float64(0, 1, 2))
	assert.Equal(2, Coalesce.Float64(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Float64(4, 1, 2, 3))
}

func TestCoalesceDuration(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, Coalesce.Duration(0, 1))
	assert.Equal(2, Coalesce.Duration(0, 1, 2))
	assert.Equal(2, Coalesce.Duration(0, 1, 2, 3))
	assert.Equal(4, Coalesce.Duration(4, 1, 2, 3))
}

func TestCoalesceTime(t *testing.T) {
	assert := assert.New(t)

	zero := time.Time{}
	one := time.Date(2018, 01, 01, 12, 00, 00, 00, time.UTC)
	two := time.Date(2018, 02, 02, 12, 00, 00, 00, time.UTC)
	three := time.Date(2018, 03, 03, 12, 00, 00, 00, time.UTC)
	four := time.Date(2018, 04, 04, 12, 00, 00, 00, time.UTC)

	assert.Equal(one, Coalesce.Time(zero, one))
	assert.Equal(two, Coalesce.Time(zero, one, two))
	assert.Equal(two, Coalesce.Time(zero, one, two, three))
	assert.Equal(four, Coalesce.Time(four, one, two, three))
}

func TestCoalesceStrings(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Coalesce.Strings(nil, nil))
	assert.NotEmpty(Coalesce.Strings(nil, []string{"foo"}))
	assert.NotEmpty(Coalesce.Strings(nil, []string{}, []string{"bar"}))
	assert.NotEmpty(Coalesce.Strings([]string{"moo"}, []string{}, []string{}))
}

func TestCoalesceBytes(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Coalesce.Bytes(nil, nil))
	assert.NotEmpty(Coalesce.Bytes(nil, []byte("foo")))
	assert.NotEmpty(Coalesce.Bytes(nil, []byte{}, []byte("bar")))
	assert.NotEmpty(Coalesce.Bytes([]byte("moo"), []byte{}, []byte{}))
}
