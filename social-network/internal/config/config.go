package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName      string `envconfig:"SERVICE_NAME" default:"social-network"`
	ServerListenPort int    `envconfig:"SERVER_LISTEN_PORT" required:"true"`
	Debug            bool   `envconfig:"DEBUG" default:"false"`

	Log Log

	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

type Log struct {
	EnableStacktrace bool `envconfig:"LOG_ENABLE_STACKTRACE" default:"false"`
}

func FromEnv() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error while parse env config: %w", err)
	}

	return cfg, nil
}
