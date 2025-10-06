package config

import "fmt"

type RabbitMQ struct {
	Address  string `envconfig:"address" required:"true"`
	Username string `envconfig:"username" required:"true"`
	Password string `envconfig:"password" required:"true"`
}

func (r *RabbitMQ) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s", r.Username, r.Password, r.Address)
}
