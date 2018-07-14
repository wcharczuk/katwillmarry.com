package request

import (
	"encoding/xml"
	"testing"

	"github.com/blend/go-sdk/assert"
)

type mockObject struct {
	XMLName      xml.Name `json:"-" xml:"Borrower"`
	ID           int      `json:"id" xml:"id,attr"`
	Email        string   `json:"email" xml:"-"`
	DeploymentID int      `json:"deployment_id" xml:"-"`
}

func testServiceRequest(t *testing.T, additionalTests ...func(*Request)) {
	assert := assert.New(t)
	sr := New().
		WithMockProvider(MockedResponseInjector).
		AsDelete().
		AsPatch().
		AsPut().
		AsPost().
		AsGet().
		WithScheme("http").
		WithHost("localhost:5001").
		WithPath("/api/v1/borrowers/2").
		WithHeader("deployment", "test").
		WithPostData("test", "regressions").
		WithQueryString("foo", "bar").
		WithTimeout(500).
		WithQueryString("moobar", "zoobar").
		WithScheme("http")

	assert.Equal("http", sr.url.Scheme)
	assert.Equal("localhost:5001", sr.url.Host)
	assert.Equal("GET", sr.method)

	req, err := sr.Request()
	assert.Nil(err)
	assert.NotNil(req)

	for _, additionalTest := range additionalTests {
		additionalTest(sr)
	}
}
