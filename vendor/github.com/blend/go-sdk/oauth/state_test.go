package oauth

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
)

func TestSerializeState(t *testing.T) {
	assert := assert.New(t)

	state := State{
		RedirectURL: "https://foo.com/bar",
		Token:       util.String.RandomLetters(32),
		SecureToken: util.String.RandomLetters(64),
	}

	contents, err := SerializeState(&state)
	assert.Nil(err)
	assert.NotEmpty(contents)

	deserialized, err := DeserializeState(contents)
	assert.Nil(err)
	assert.NotNil(deserialized)
	assert.Equal(state.RedirectURL, deserialized.RedirectURL)
	assert.Equal(state.Token, deserialized.Token)
	assert.Equal(state.SecureToken, deserialized.SecureToken)
}
