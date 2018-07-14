package web

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestViewCacheAddRawViews(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	vc.AddLiterals(`{{ define "test" }}<h1> This is a test. </h1>{{ end }}`)

	view, err := vc.Parse()
	assert.Nil(err)
	assert.NotNil(view)

	buf := bytes.NewBuffer(nil)
	assert.Nil(view.ExecuteTemplate(buf, "test", nil))
	assert.NotEmpty(buf.String())
}

func TestViewCacheCached(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.True(vc.Cached())
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.viewLiterals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())
}

func TestViewCacheCachingDisabled(t *testing.T) {
	assert := assert.New(t)

	vc := NewViewCache()
	assert.True(vc.Cached())
	vc.SetCached(false)
	vc.AddLiterals(`{{ define "foo" }}bar{{ end }}`)
	assert.Nil(vc.Initialize())

	tmp, err := vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf := bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("bar", buf.String())

	vc.viewLiterals = []string{`{{ define "foo" }}baz{{ end }}`}
	tmp, err = vc.Lookup("foo")
	assert.Nil(err)
	assert.NotNil(tmp)
	buf = bytes.NewBuffer(nil)
	assert.Nil(tmp.Execute(buf, nil))
	assert.Equal("baz", buf.String())
}
