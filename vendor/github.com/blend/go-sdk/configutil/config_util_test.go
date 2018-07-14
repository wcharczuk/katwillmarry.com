package configutil

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/uuid"
)

type config struct {
	Environment string `json:"env" yaml:"env" env:"SERVICE_ENV"`
	Other       string `json:"other" yaml:"other" env:"OTHER"`
}

func TestPath(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Delete(EnvVarConfigPath)

	assert.Empty(Path())
	assert.Equal("testdata/foo.yml", Path("testdata/foo.yml"))
}

func TestPathEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()
	env.Env().Set(EnvVarConfigPath, "testdata/config.yml")
	assert.Equal("testdata/config.yml", Path())
}

func TestDeserializeInvalid(t *testing.T) {
	assert := assert.New(t)

	err := Deserialize(".???", nil, nil)
	assert.NotNil(err)
	assert.True(IsInvalidConfigExtension(err))
}

func TestDeserializeYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(Deserialize(ExtensionYAML, bytes.NewBuffer([]byte("env: test\nother: foo")), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestDeserializeYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(Deserialize(ExtensionYML, bytes.NewBuffer([]byte("env: test\nother: foo")), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestDeserializeJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	assert.Nil(Deserialize(ExtensionJSON, bytes.NewBuffer([]byte(`{"env": "test", "other": "foo"}`)), &cfg))
	assert.Equal("test", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestRead(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set(env.VarServiceEnv, "dev")
	var cfg config
	err := ReadFromReader(&cfg, bytes.NewBuffer([]byte("env: test\nother: foo")), ExtensionYAML)
	assert.Nil(err)
	assert.Equal("dev", cfg.Environment)
}

func TestReadFromPathYAML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	err := ReadFromPath(&cfg, "testdata/config.yaml")
	assert.Nil(err)
	assert.Equal("test_yaml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestReadFromPathYML(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	err := ReadFromPath(&cfg, "testdata/config.yml")
	assert.Nil(err)
	assert.Equal("test_yml", cfg.Environment)
	assert.Equal("foo", cfg.Other)
}

func TestReadFromPathJSON(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	err := ReadFromPath(&cfg, "testdata/config.json")
	assert.Nil(err)
	assert.Equal("test_json", cfg.Environment)
	assert.Equal("moo", cfg.Other)
}

func TestReadPathUnset(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set(env.VarServiceEnv, "dev")
	var cfg config
	err := ReadFromPath(&cfg, "")
	assert.True(IsConfigPathUnset(err))
	assert.Equal("dev", cfg.Environment)
}

func TestReadPathNotFound(t *testing.T) {
	assert := assert.New(t)

	var cfg config
	err := ReadFromPath(&cfg, filepath.Join("testdata", uuid.V4().String()))
	assert.True(IsNotExist(err))
}

func TestIsUnset(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsConfigPathUnset(exception.New(ErrConfigPathUnset)))
	assert.False(IsConfigPathUnset(exception.New(uuid.V4().String())))
}

func TestIsIgnored(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsIgnored(nil))
	assert.True(IsIgnored(exception.New(ErrConfigPathUnset)))
	assert.True(IsIgnored(exception.New(ErrInvalidConfigExtension)))
}
