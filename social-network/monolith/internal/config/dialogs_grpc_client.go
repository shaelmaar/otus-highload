package config

import "time"

type DialogsGRPCClient struct {
	Host    string        `envconfig:"HOST" required:"true"`
	TLS     bool          `envconfig:"TLS" default:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"5s"`
}
