package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("AUTH_SECRET", Base64Encode(util.Crypto.MustCreateKey(32)))

	config := NewConfigFromEnv()
	assert.NotEmpty(config.GetAuthSecret())
}
