package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLogWriterWrite(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = false
	writer.useColor = false
	writer.Write(Messagef(Info, "test string"))
	assert.Equal("[info] test string\n", string(buffer.Bytes()))
}

func TestLogWriterWriteWithLabel(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = true
	writer.useColor = false
	writer.Write(Messagef(Info, "test string").WithHeadings("unit-test"))
	assert.Equal("[unit-test] [info] test string\n", string(buffer.Bytes()))
}

func TestLogWriterWriteWithLabelColorized(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.showHeadings = true
	writer.useColor = true
	writer.Write(Messagef(Info, "test string").WithHeadings("unit-test"))
	assert.Equal("["+ColorBlue.Apply("unit-test")+"] ["+ColorLightWhite.Apply("info")+"] test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutputCoalesced(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.useColor = false
	writer.WriteError(Messagef(Error, "test %s", "string"))
	assert.Equal("[error] test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutput(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	writer := NewTextWriter(stdout).WithErrorOutput(stderr)
	writer.showTimestamp = false
	writer.useColor = false

	writer.WriteError(Messagef(Error, "test %s", "string"))
	assert.Equal(0, stdout.Len())
	assert.Equal("[error] test string\n", string(stderr.Bytes()))
}

func TestWriterLabels(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewTextWriter(buffer)
	writer.showTimestamp = false
	writer.useColor = false
	writer.WriteError(Messagef(Error, "test %s", "string").WithLabel("foo", "bar").WithLabel("moo", "boo"))
	assert.True(strings.HasPrefix(buffer.String(), "[error] test string"))
}
