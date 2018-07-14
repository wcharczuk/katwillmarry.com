package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHTTPSUpgrader(t *testing.T) {
	assert := assert.New(t)

	upgrader := NewHTTPSUpgrader()

	ts := httptest.NewServer(upgrader)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	assert.NotNil(err)
}
