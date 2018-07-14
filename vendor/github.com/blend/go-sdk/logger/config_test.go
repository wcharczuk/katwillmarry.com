package logger

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/yaml"
)

func b(v bool) *bool {
	return &v
}

func TestConfigProperties(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(Config{}.GetHeading())
	assert.Equal("test", Config{Heading: "test"}.GetHeading())

	assert.Equal(OutputFormatText, Config{}.GetOutputFormat())
	assert.Equal(OutputFormatJSON, Config{OutputFormat: string(OutputFormatJSON)}.GetOutputFormat())

	assert.Equal(AsStrings(DefaultFlags...), Config{}.GetFlags())
	assert.Equal([]string{"foo", "bar"}, Config{Flags: []string{"foo", "bar"}}.GetFlags())

	assert.Equal(DefaultRecoverPanics, Config{}.GetRecoverPanics())
	assert.Equal(!DefaultRecoverPanics, Config{}.GetRecoverPanics(!DefaultRecoverPanics))
	assert.Equal(!DefaultRecoverPanics, Config{RecoverPanics: b(!DefaultRecoverPanics)}.GetRecoverPanics())

	assert.Equal(DefaultWorkerQueueDepth, Config{}.GetQueueDepth())
	assert.Equal(DefaultWorkerQueueDepth>>1, Config{}.GetQueueDepth(DefaultWorkerQueueDepth>>1))
	assert.Equal(DefaultWorkerQueueDepth>>2, Config{QueueDepth: DefaultWorkerQueueDepth >> 2}.GetQueueDepth(DefaultWorkerQueueDepth>>1))
}

func TestNewConfigFlags(t *testing.T) {
	assert := assert.New(t)

	cfg := &Config{Flags: []string{"all", "-debug"}}
	log := NewFromConfig(cfg)
	defer log.Close()

	assert.True(log.IsEnabled(Silly))
	assert.True(log.IsEnabled(Info))
	assert.False(log.IsEnabled(Debug))
	assert.True(log.IsEnabled(Warning))
	assert.True(log.IsEnabled(Error))
	assert.True(log.IsEnabled(Fatal))
	assert.True(log.IsEnabled(Audit))
	assert.True(log.IsEnabled(HTTPRequest))
}

func TestConfigYAML(t *testing.T) {
	assert := assert.New(t)

	corpus := `
heading: test-heading
outputFormat: test-format
flags: [ "foo", "bar" ]
hiddenFlags: [ "buzz", "wuzz" ]
recoverPanics: false
queueDepth: 256
`

	var cfg Config
	assert.Nil(yaml.Unmarshal([]byte(corpus), &cfg))
	assert.Equal("test-heading", cfg.GetHeading())
	assert.Equal("test-format", cfg.GetOutputFormat())
	assert.Equal([]string{"foo", "bar"}, cfg.GetFlags())
	assert.Equal([]string{"buzz", "wuzz"}, cfg.GetHiddenFlags())
	assert.False(cfg.GetRecoverPanics())
	assert.Equal(256, cfg.GetQueueDepth())
}

func TestConfigJSON(t *testing.T) {
	assert := assert.New(t)

	corpus := `{
"heading": "test-heading",
"outputFormat": "test-format",
"flags": [ "foo", "bar" ],
"hiddenFlags": [ "buzz", "wuzz" ],
"recoverPanics": false,
"queueDepth": 256
}
`

	var cfg Config
	assert.Nil(json.Unmarshal([]byte(corpus), &cfg))
	assert.Equal("test-heading", cfg.GetHeading())
	assert.Equal("test-format", cfg.GetOutputFormat())
	assert.Equal([]string{"foo", "bar"}, cfg.GetFlags())
	assert.Equal([]string{"buzz", "wuzz"}, cfg.GetHiddenFlags())
	assert.False(cfg.GetRecoverPanics())
	assert.Equal(256, cfg.GetQueueDepth())
}

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)

	cfg := NewConfigFromEnv()

	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "http.response"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "info"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "warning"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "error"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "fatal"
	})
	assert.None(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "debug"
	})
}

func TestNewConfigFromEnvWithVars(t *testing.T) {
	assert := assert.New(t)

	env.SetEnv(env.Vars{
		"LOG_EVENTS": "info,debug,error,test",
		"LOG_HIDDEN": "debug",
	})
	defer env.Restore()

	cfg := NewConfigFromEnv()

	assert.NotEmpty(cfg.GetFlags())
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "info"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "debug"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "error"
	})
	assert.Any(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "test"
	})
	assert.None(cfg.GetFlags(), func(v Any) bool {
		return v.(string) == "fatal"
	})

}

func TestGetWritersWithOutputFormat(t *testing.T) {
	assert := assert.New(t)

	config := &Config{OutputFormat: string(OutputFormatJSON)}
	writers := config.GetWriters()
	assert.Len(writers, 1)
	assert.Equal(OutputFormatJSON, writers[0].OutputFormat())
	config.OutputFormat = string(OutputFormatText)
	writers = config.GetWriters()
	assert.Len(writers, 1)
	assert.Equal(OutputFormatText, writers[0].OutputFormat())
	config.OutputFormat = "nope"
	writers = config.GetWriters()
	assert.Len(writers, 1)
	assert.Equal(OutputFormatText, writers[0].OutputFormat())
	config.OutputFormat = strings.ToUpper(string(OutputFormatJSON))
	writers = config.GetWriters()
	assert.Len(writers, 1)
	assert.Equal(OutputFormatJSON, writers[0].OutputFormat())
}

func TestNewJSONWriterConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("LOG_JSON_PRETTY", "false")
	cfg := NewJSONWriterConfigFromEnv()
	assert.False(cfg.GetPretty())
}
