package logger

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestNewEventMeta(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)
	assert.Equal(Info, em.Flag())
	assert.False(em.Timestamp().IsZero())
}

func TestEventMetaProperties(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)

	assert.NotNil(em.Annotations())
	assert.NotNil(em.Labels())

	assert.Empty(em.Headings())
	em.SetHeadings("Headings")
	assert.Equal([]string{"Headings"}, em.Headings())

	assert.Equal(Info, em.Flag())
	em.SetFlag(Fatal)
	assert.Equal(Fatal, em.Flag())

	assert.Empty(em.FlagTextColor())
	em.SetFlagTextColor(ColorRed)
	assert.Equal(ColorRed, em.FlagTextColor())

	assert.False(em.Timestamp().IsZero())
	em.SetTimestamp(time.Time{})
	assert.True(em.Timestamp().IsZero())

	assert.Empty(em.Labels())
	em.AddLabelValue("foo", "bar")
	assert.Equal("bar", em.Labels()["foo"])

	em.SetLabels(nil)
	assert.Empty(em.Labels())

	assert.Empty(em.Annotations())
	em.AddAnnotationValue("buzz", "fuzz")
	assert.Equal("fuzz", em.Annotations()["buzz"])

	em.SetAnnotations(nil)
	assert.Empty(em.Annotations())
}
