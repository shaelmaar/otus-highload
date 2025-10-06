package config

import "time"

type Valkey struct {
	Address    string        `envconfig:"ADDRESS"     required:"true"`
	DB         int           `envconfig:"DB"          required:"true"`
	SetTimeout time.Duration `envconfig:"SET_TIMEOUT" required:"true"`
	Password   string        `envconfig:"PASSWORD" required:"true"`
}
