package logger

import (
	"bytes"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(nil, func(e Event) {
		defer wg.Done()
		didFire = true

		typed, isTyped := e.(*MessageEvent)
		assert.True(isTyped)
		assert.Equal("test", typed.Message())
	}, DefaultWorkerQueueDepth)

	w.Start()
	defer w.Close()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
}

func TestWorkerPanics(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	wr := NewTextWriter(buffer)

	log := New().WithFlags(AllFlags()).WithWriter(wr)
	defer log.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(log, func(e Event) {
		defer wg.Done()
		didFire = true
		panic("only a test")
	}, DefaultWorkerQueueDepth)
	w.Start()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
	w.Close()
	assert.NotEmpty(buffer.String())
}
