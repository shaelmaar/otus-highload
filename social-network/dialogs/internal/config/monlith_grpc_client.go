package config

import "time"

type MonolithGRPCClient struct {
	Host    string        `envconfig:"HOST" required:"true"`
	TLS     bool          `envconfig:"TLS" default:"true"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"5s"`
}
