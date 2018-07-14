package web

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTLSConfigGetConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := TLSConfig{
		CertPath: "testdata/testcert.pem",
		KeyPath:  "testdata/testkey.pem",
	}

	tlsConfig, err := cfg.GetConfig()
	assert.Nil(err)
	assert.NotNil(tlsConfig)
	assert.NotEmpty(tlsConfig.Certificates)
}
