package util

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimeUnixMillis(t *testing.T) {
	assert := assert.New(t)

	zero := Time.UnixMillis(time.Unix(0, 0))
	assert.Zero(zero)

	sample := Time.UnixMillis(time.Date(2015, 03, 07, 16, 0, 0, 0, time.UTC))
	assert.Equal(1425744000000, sample)
}

func TestTimeFromUnixMillis(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2015, 03, 07, 16, 0, 0, 0, time.UTC)
	millis := Time.UnixMillis(ts)
	ts2 := Time.FromUnixMillis(millis)
	assert.Equal(ts, ts2)
}

func TestTimeIsWeekday(t *testing.T) {
	assert := assert.New(t)

	monday := time.Date(2018, 05, 21, 12, 0, 0, 0, time.UTC)
	tuesday := time.Date(2018, 05, 22, 12, 0, 0, 0, time.UTC)
	wednesday := time.Date(2018, 05, 23, 12, 0, 0, 0, time.UTC)
	thursday := time.Date(2018, 05, 24, 12, 0, 0, 0, time.UTC)
	friday := time.Date(2018, 05, 25, 12, 0, 0, 0, time.UTC)
	saturday := time.Date(2018, 05, 26, 12, 0, 0, 0, time.UTC)
	sunday := time.Date(2018, 05, 27, 12, 0, 0, 0, time.UTC)
	assert.True(Time.IsWeekDay(monday.Weekday()))
	assert.True(Time.IsWeekDay(tuesday.Weekday()))
	assert.True(Time.IsWeekDay(wednesday.Weekday()))
	assert.True(Time.IsWeekDay(thursday.Weekday()))
	assert.True(Time.IsWeekDay(friday.Weekday()))
	assert.False(Time.IsWeekDay(saturday.Weekday()))
	assert.True(Time.IsWeekendDay(saturday.Weekday()))
	assert.False(Time.IsWeekDay(sunday.Weekday()))
	assert.True(Time.IsWeekendDay(sunday.Weekday()))
}
