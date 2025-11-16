package config

import (
	"fmt"
	"math/rand"
)

type RabbitMQ struct {
	Addresses []string `envconfig:"address" required:"true"`
	Username  string   `envconfig:"username" required:"true"`
	Password  string   `envconfig:"password" required:"true"`
}

func (r RabbitMQ) Validate() error {
	if len(r.Addresses) == 0 {
		return fmt.Errorf("rabbitMQ addresses are empty")
	}

	return nil
}

func (r *RabbitMQ) URL() string {
	i := rand.Intn(len(r.Addresses)) //nolint:gosec // вместо балансировки просто выбор случайного.

	return fmt.Sprintf("amqp://%s:%s@%s", r.Username, r.Password, r.Addresses[i])
}
