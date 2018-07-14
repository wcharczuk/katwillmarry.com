package request

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//--------------------------------------------------------------------------------
// RequestMeta
//--------------------------------------------------------------------------------

// NewRequestMeta returns a new meta object for a request.
func NewRequestMeta(req *http.Request) *Meta {
	return &Meta{
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Header,
	}
}

// NewRequestMetaWithBody returns a new meta object for a request and reads the body.
func NewRequestMetaWithBody(req *http.Request) (*Meta, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	return &Meta{
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Header,
		Body:    body,
	}, nil
}

// Meta is a summary of the request meta useful for logging.
type Meta struct {
	StartTime time.Time
	Method    string
	URL       *url.URL
	Headers   http.Header
	Body      []byte
}
