package config

import "time"

type ValkeyDB struct {
	Addresses     []string      `envconfig:"ADDRESSES" required:"true"`
	Password      string        `envconfig:"PASSWORD"`
	WriteTimeout  time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	SlavePoolSize int           `envconfig:"SLAVE_POOL_SIZE" default:"10"`
}
