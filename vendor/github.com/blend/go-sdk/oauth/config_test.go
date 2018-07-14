package oauth

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("OAUTH_REDIRECT_URI", "https://app.com/oauth/google")
	env.Env().Set("OAUTH_HOSTED_DOMAIN", "foo.com")
	env.Env().Set("OAUTH_CLIENT_ID", "foo")
	env.Env().Set("OAUTH_CLIENT_SECRET", "bar")

	cfg := NewConfigFromEnv()
	assert.Equal("foo", cfg.GetClientID())
	assert.Equal("bar", cfg.GetClientSecret())
	assert.Equal("https://app.com/oauth/google", cfg.GetRedirectURI())
	assert.Equal("foo.com", cfg.GetHostedDomain())
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.True(Config{ClientID: "foo"}.IsZero())
	assert.False(Config{ClientID: "foo", ClientSecret: "bar"}.IsZero())
}
