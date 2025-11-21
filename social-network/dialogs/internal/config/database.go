package config

type TarantoolDB struct {
	Host string `envconfig:"HOST" required:"true"`
	Port string `envconfig:"PORT" required:"true"`
	User string `envconfig:"USER" required:"true"`
	Pass string `envconfig:"PASS"`
}
