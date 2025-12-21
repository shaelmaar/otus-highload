package config

type Kafka struct {
	Brokers   []string `envconfig:"BROKERS" required:"true"`
	GroupName string   `envconfig:"GROUPNAME" default:"dialogs"`
}
