package rabbitmq

import "time"

type ProducerOption[T Message] func(producer *Producer[T])

func WithMessageTTL[T Message](ttl time.Duration) ProducerOption[T] {
	return func(producer *Producer[T]) {
		producer.messageTTL = ttl
	}
}
