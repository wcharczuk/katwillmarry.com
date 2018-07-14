package config

import (
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/web"
)

type Config struct {
	Web    web.Config    `yaml:"web"`
	DB     db.Config     `yaml:"db"`
	OAuth  oauth.Config  `yaml:"oauth"`
	Logger logger.Config `yaml:"logger"`
}
