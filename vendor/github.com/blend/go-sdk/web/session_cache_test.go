package web

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
)

func TestSessionCacheAdd(t *testing.T) {
	assert := assert.New(t)

	cache := NewSessionCache()
	assert.False(cache.IsActive(uuid.V4().String()))
	cache.Upsert(&Session{
		SessionID: "foo",
	})
	assert.True(cache.IsActive("foo"))
}

func TestSessionCacheGet(t *testing.T) {
	assert := assert.New(t)

	cache := NewSessionCache()
	assert.False(cache.IsActive(uuid.V4().String()))
	cache.Upsert(&Session{
		SessionID: "foo",
	})
	got := cache.Get("foo")
	assert.NotNil(got)
	assert.Equal("foo", got.SessionID)

	got = cache.Get(uuid.V4().String())
	assert.Nil(got)
}

func TestSessionCacheRemove(t *testing.T) {
	assert := assert.New(t)

	cache := NewSessionCache()
	cache.Upsert(&Session{
		SessionID: "foo",
	})
	assert.True(cache.IsActive("foo"))
	cache.Remove("foo")
	assert.False(cache.IsActive("foo"))
}

func TestSessionCacheIsActive(t *testing.T) {
	assert := assert.New(t)

	cache := NewSessionCache()

	assert.False(cache.IsActive(uuid.V4().String()))

	cache.Upsert(&Session{
		SessionID: "foo",
	})
	assert.True(cache.IsActive("foo"))

	cache.Upsert(&Session{
		SessionID:  "bar",
		ExpiresUTC: util.OptionalTime(time.Now().UTC().Add(time.Hour)),
	})
	assert.True(cache.IsActive("bar"))

	cache.Upsert(&Session{
		SessionID:  "baz",
		ExpiresUTC: util.OptionalTime(time.Now().UTC().Add(-time.Hour)),
	})
	assert.False(cache.IsActive("baz"))
}
